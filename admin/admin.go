package admin

import (
	"net/http"
	"html/template"
	"onlinestore/auth"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/smtp"
	"os"
	"github.com/joho/godotenv"
	"fmt"
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
		r.ParseForm()
		users := GetUsers(client, database, collection, bson.M{}, bson.D{})
		for _, user := range users {
			err := sendEmail(r.FormValue("subject"), r.FormValue("message"), user.Email)
			if err != nil {
				http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
				return
			}

		}
		http.Redirect(w, r, "http://127.0.0.1:8080/admin", http.StatusSeeOther)
		return
	}
}

func GetUsers(client *mongo.Client, database, collection string, filter bson.M, sorting bson.D) []auth.User {
    coll := client.Database(database).Collection(collection)

    findOptions := options.Find().SetSort(sorting)

    cursor, err := coll.Find(context.TODO(), filter, findOptions)
    if err != nil {
        panic(err)
    }

    var users []auth.User
    if err := cursor.All(context.TODO(), &users); err != nil {
        panic(err)
    }

    return users
}

func sendEmail(subject, message string, recipient string) error {
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
