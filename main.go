package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"onlinestore/db"
	"onlinestore/products"
	"github.com/gorilla/mux"
	"onlinestore/auth"
	"os"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"go.mongodb.org/mongo-driver/mongo"
	"onlinestore/admin"
)

type GetMessage struct {
    Message	string
}

type PostMessage struct {
	Status	string
	Message string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getMessage := GetMessage{
			Message: "Hello, server! This is JSON data from Postman.",
		}
		jsonGetMessage, err := json.Marshal(getMessage)
		if err != nil {
			getMessage.Message = "Invalid JSON message"
			jsonGetMessage, _ = json.Marshal(getMessage)
		}
		w.Write(jsonGetMessage)

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		postValue := r.FormValue("message")
		postMessage := PostMessage {
			Status: "fail",
			Message: "Invalid JSON message",
		}
		if postValue != "" {
			_, err := strconv.Atoi(postValue)
			if err != nil {
				postMessage.Status = "success"
				postMessage.Message ="Data successfully received"
			}
		}
		jsonGetMessage, _ := json.Marshal(postMessage)
		w.Write(jsonGetMessage)
	}
}


var limiter = rate.NewLimiter(1, 1)
func productsLimiter(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string, logs *logrus.Logger) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	ProtectedHandler(w, r)
	products.ProductsHandler(w, r, client, database, collection, logs)
}

func main() {
	logs := logrus.New()
	logs.SetFormatter(&logrus.JSONFormatter{})
	
	// Log output to a file
	file, err := os.OpenFile("user_actions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logs.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	logs.SetOutput(file)
	uri := "mongodb://localhost:27017"
	client, dbErr := db.ConnectMongoDB(uri)
	if dbErr != nil {
		log.Println("Could not connect to the MongoDB")
	}
	database := "onlineStore"
	collection := "AlisherExpress"
	mux := mux.NewRouter()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		productsLimiter(w, r, client, database, collection, logs)
	})

	mux.HandleFunc("/products/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
        products.Product(w, r, client, database, collection, logs)
    })

	admin.RegisterRoutes(mux, client, database, "Users")
	auth.RegisterRoutes(mux, client, database, "Users")
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}


func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "http://127.0.0.1:8080/login", http.StatusSeeOther)
		return
	}

	err = auth.VerifyToken(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "http://127.0.0.1:8080/login", http.StatusSeeOther)
		return
	}

}

