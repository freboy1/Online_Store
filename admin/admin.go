package admin

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"onlinestore/models"
	"onlinestore/products"
	"os"
	"strconv"
	"time"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/admin.html")
		tmpl.Execute(w, map[string]interface{}{})
	}
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/admin-send.html")
		tmpl.Execute(w, map[string]interface{}{})
	} else if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // Limit file size to 10 MB
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("attachment")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Error processing file", http.StatusBadRequest)
			return
		}
		var photoData []byte
		if file != nil {
			defer file.Close()
			photoData, err = io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}
		}


		users := GetUsers(client, database, collection, bson.M{}, bson.D{})
		for _, user := range users {
			emailErr := sendEmailImage(
				r.FormValue("subject"),    // Subject from the form
				r.FormValue("message"),    // Message from the form
				user.Email,                // User email from database
				header.Filename,           // Attachment filename
				photoData,                 // Photo data
			)
			if emailErr != nil {
				http.Error(w, "Failed to send email: "+emailErr.Error(), http.StatusInternalServerError)
				return
			}

		}
		http.Redirect(w, r, "http://127.0.0.1:8080/admin", http.StatusSeeOther)
		return
	}
}

func GetUsers(client *mongo.Client, database, collection string, filter bson.M, sorting bson.D) []models.User {
    coll := client.Database(database).Collection(collection)

    findOptions := options.Find().SetSort(sorting)

    cursor, err := coll.Find(context.TODO(), filter, findOptions)
    if err != nil {
        panic(err)
    }

    var users []models.User
    if err := cursor.All(context.TODO(), &users); err != nil {
        panic(err)
    }

    return users
}

func CreateUser(client *mongo.Client, ctx context.Context, dataBase, col string, user models.User) (*mongo.InsertOneResult, error) {
	collection := client.Database(dataBase).Collection(col)
	user.Id = uuid.New() 
    result, err := collection.InsertOne(ctx, user)
    return result, err
}

func SendEmail(subject, message string, recipient string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, message))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, msg)
}


func sendEmailImage(subject, message, recipient, filename string, photoData []byte) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	encodedPhoto := base64.StdEncoding.EncodeToString(photoData)

    // Create the MIME email
    email := fmt.Sprintf(
        "From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=boundary\n\n"+
            "--boundary\nContent-Type: text/plain; charset=utf-8\n\n%s\n\n"+
            "--boundary\nContent-Type: image/jpeg\nContent-Transfer-Encoding: base64\nContent-Disposition: attachment; filename=\"%s\"\n\n%s\n--boundary--",
        from, recipient, subject, message, filename, encodedPhoto,
    )

    // Send the email
    auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, []byte(email))

}

func UsersHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	users := GetUsers(client, database, collection, bson.M{}, bson.D{})
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/admin-users.html")
		tmpl.Execute(w, users)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		name, email, password, cashStr, role := r.FormValue("name"), r.FormValue("email"), r.FormValue("password"), r.FormValue("cash"), r.FormValue("role")
		cash, _ := strconv.Atoi(cashStr)
		user := models.User{
			Username: name,
			Email: email,
			Password: password,
			Cash: cash,
			Role: role,
		}
		result, err := insertUser(client, context.TODO(), database, collection, user)
		if err != nil {
			tmpl, _ := template.ParseFiles("templates/admin-users.html")
			tmpl.Execute(w, users)
			return

		}
	
		fmt.Println("Inserted user with ID:", result.InsertedID)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
}

func insertUser (client *mongo.Client, ctx context.Context, dataBase, col string, user models.User) (*mongo.InsertOneResult, error) {

    // select database and collection ith Client.Database method 
    // and Database.Collection method
    collection := client.Database(dataBase).Collection(col)
	user.Id = uuid.New() 
	user.Code = generateRandomCode(4)
	user.Verified = "false"
	user.Products = make([]products.ProductModel, 0)
    // InsertOne accept two argument of type Context 
    // and of empty interface   
    result, err := collection.InsertOne(ctx, user)
    return result, err
}

func generateRandomCode(length int) string {
	if length <= 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10)) // Append a random digit (0-9)
	}
	return code
}

func UserHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	vars := mux.Vars(r)
	id := vars["id"] 
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		fmt.Printf("Error parsing UUID: %v\n", err)
		return
	}
	mongoUUID := primitive.Binary{
		Subtype: 0x00,       // MongoDB UUID subtype
		Data:    parsedUUID[:], // UUID as byte array
	}

	users := GetUsers(client, database, collection, bson.M{"id": mongoUUID}, bson.D{})
	tmpl, _ := template.ParseFiles("templates/admin-user.html")
	if len(users) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user := users[0]
	if r.Method == http.MethodGet {
		tmpl.Execute(w, user)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			tmpl.Execute(w, user)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "delete":
			result, err := deleteOne(client, context.TODO(), database, collection, mongoUUID)
			if err != nil {
				tmpl, _ := template.ParseFiles("templates/admin-user.html")
				tmpl.Execute(w, user)
				return
			}
		
			fmt.Println("Deleted succesfully:", result)
			http.Redirect(w, r, "http://127.0.0.1:8080/admin/users", http.StatusSeeOther)
			return
		case "update":
			if err := r.ParseForm(); err != nil {
				tmpl.Execute(w, user)
				return
			}
			name, email, cashStr, role := r.FormValue("name"), r.FormValue("email"), r.FormValue("cash"), r.FormValue("role")
			cash, _ := strconv.Atoi(cashStr)
			userNew := models.User{
				Username: name,
				Email: email,
				Cash: cash,
				Role: role,
			}
			err = updateOne(client, context.TODO(), database, collection, mongoUUID, userNew)
			if err != nil {
				tmpl, _ := template.ParseFiles("templates/admin-user.html")
				tmpl.Execute(w, user)
				return
			}
			
			http.Redirect(w, r, "http://127.0.0.1:8080/admin/users", http.StatusSeeOther)
			return
		}
	}
}

func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col string, id primitive.Binary) (*mongo.DeleteResult, error) {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"id", id}}
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}

func updateOne(client *mongo.Client, ctx context.Context, dataBase, col string, id primitive.Binary, User models.User) error {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"username", User.Username},
			{"email", User.Email},
			{"cash", User.Cash},
			{"role", User.Role},
		}},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("failed to update product")
		return err
	}

	// Check if the product was found and updated
	if result.MatchedCount == 0 {
		return err
	}

	fmt.Printf("Successfully updated %d product(s)\n", result.ModifiedCount)
	return nil
}