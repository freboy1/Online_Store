package auth

import (
	"html/template"
	"net/http"
	"onlinestore/admin"
	"onlinestore/models"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")
	user, err := ExistUser(client, database, collection, email, password)
    if  err == nil{
        tokenString, err := CreateToken(user.Email, user.Role)
        if err != nil {
            http.Error(w, "Error creating token", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     "auth_token",
            Value:    tokenString,
            HttpOnly: true,
            Secure:   false,
            Path:     "/",
            Expires:  time.Now().Add(24 * time.Hour),
        })
        http.Redirect(w, r, "http://127.0.0.1:8080/products", http.StatusSeeOther)
		return
    }
	http.Redirect(w, r, "http://127.0.0.1:8080/login", http.StatusSeeOther)
}


func LoginGet(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, map[string]interface{}{})
}


func ExistUser(client *mongo.Client, database, collection, email, password string) (models.User, error) {
	user := admin.GetUsers(client, database, collection, bson.M{"email": email, "password": password}, bson.D{})
	if len(user) != 0 {
		return user[0], nil
	}
	return models.User{}, fmt.Errorf("NO user")
}