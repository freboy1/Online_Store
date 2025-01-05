package auth

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"onlinestore/admin"
	"onlinestore/models"
	"onlinestore/products"
	"time"

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
	if err == nil {
		tokenString, err := CreateToken(user.Email, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}

		SetAuthCookie(w, tokenString)
		http.Redirect(w, r, "http://127.0.0.1:8080/products", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "http://127.0.0.1:8080/login", http.StatusSeeOther)
}

func SetAuthCookie(w http.ResponseWriter, tokenString string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
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

func RegisterHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, map[string]interface{}{})
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username, email, password := r.FormValue("username"), r.FormValue("email"), r.FormValue("password")
		user := models.User{
			Username: username,
			Email:    email,
			Password: password,
			Role:     "user",
			Code:     GenerateRandomCode(4),
			Verified: "no",
			Cash:     0,
			Products: make([]products.ProductModel, 0),
		}
		admin.CreateUser(client, context.TODO(), database, collection, user)
		tokenString, err := CreateToken(user.Email, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}

		SetAuthCookie(w, tokenString)
		http.Redirect(w, r, "http://127.0.0.1:8080/register/email-verification", http.StatusSeeOther)
		return
	}
}

func EmailVerification(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	cookie, err := r.Cookie("auth_token")
	email, _ := GetClaim(cookie.Value, "email")
	password, _ := GetClaim(cookie.Value, "password")
	user, _ := ExistUser(client, database, collection, email, password)
	if r.Method == http.MethodGet {
		admin.SendEmail("Your Code", "Hello Here is your code for verification" + user.Code, user.Email)
		tmpl, _ := template.ParseFiles("templates/email-verification.html")
		tmpl.Execute(w, map[string]interface{}{})
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		code := r.FormValue("code")
		if err != nil {
			http.Redirect(w, r, "http://127.0.0.1:8080/register", http.StatusSeeOther)
			return
		}
		if user.Code == code {
			http.Redirect(w, r, "http://127.0.0.1:8080/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "http://127.0.0.1:8080/register/email-verification", http.StatusSeeOther)
		return
	}
}

func GenerateRandomCode(length int) string {
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
