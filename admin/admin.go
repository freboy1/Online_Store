package admin

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"onlinestore/models"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/google/uuid"
	"io"
	"encoding/base64"
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
	if r.Method == http.MethodGet {
		users := GetUsers(client, database, collection, bson.M{}, bson.D{})
		tmpl, _ := template.ParseFiles("templates/admin-users.html")
		tmpl.Execute(w, users)
	}
}