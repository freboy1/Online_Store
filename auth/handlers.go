package auth

import (
	"fmt"
	"net/http"
	"time"
	"html/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"onlinestore/admin"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    if ExistUser(client, database, collection, email, password) {
        tokenString, err := CreateToken(email)
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



func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized: Missing or invalid token")
		return
	}

	err = verifyToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized: Invalid token")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to the protected area")
}

func ExistUser(client *mongo.Client, database, collection, email, password string) bool {
	user := admin.GetUsers(client, database, collection, bson.M{"email": email, "password": password}, bson.D{})
	return len(user) != 0
}