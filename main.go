package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"html/template"
	"log"
	"net/http"
	"onlinestore/admin"
	"onlinestore/auth"
	"onlinestore/cart"
	"onlinestore/chat"
	"onlinestore/db"
	"onlinestore/products"
	"os"
	"strconv"
	"time"
)

type GetMessage struct {
	Message string
}

type PostMessage struct {
	Status  string
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
		postMessage := PostMessage{
			Status:  "fail",
			Message: "Invalid JSON message",
		}
		if postValue != "" {
			_, err := strconv.Atoi(postValue)
			if err != nil {
				postMessage.Status = "success"
				postMessage.Message = "Data successfully received"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Error loading .env file")
		return
	}
	// Log output to a file
	file, err := os.OpenFile("user_actions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logs.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	logs.SetOutput(file)
	uri := os.Getenv("MONGO_URI")
	client, dbErr := db.ConnectMongoDB(uri)
	if dbErr != nil {
		log.Println("Could not connect to the MongoDB")
	} else {
		log.Println("Connected to the MongoDB")
	}
	database := "onlineStore"
	collection := "AlisherExpress"

	mux := mux.NewRouter()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/ws", chat.HandleConnections)
	mux.HandleFunc("/messageshistory", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleMessageHistory(w, r, client, database, "Chats")
	})
	mux.HandleFunc("/messageshistory/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleGetMessageHistory(w, r, client, database, "Chats")
	})
	mux.HandleFunc("/chat/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("templates/chat.html")
		temp.Execute(w, nil)
	})

	mux.HandleFunc("/admin/chats", func(w http.ResponseWriter, r *http.Request) {
		chats := chat.GetChats(client, database, "Chats", bson.M{"status": "Active"}, bson.D{})
		temp, _ := template.ParseFiles("templates/chats.html")
		temp.Execute(w, chats)
	})
	mux.HandleFunc("/getrole", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("auth_token")
		role, _ := admin.GetClaim(cookie.Value, "role")
		w.Header().Set("Content-Type", "applications/json")
		json.NewEncoder(w).Encode(map[string]string{"role": role})

	})
	mux.HandleFunc("/closechat", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		chatID, _ := strconv.ParseInt(r.FormValue("chat_id"), 10, 64)
		dataBase := client.Database(database).Collection("Chats")
		dataBase.UpdateOne(r.Context(), bson.M{"id": chatID}, bson.M{"$set": bson.M{"status": "Closed"}})
		http.Redirect(w, r, "/admin/chats", http.StatusSeeOther)
	})

	mux.HandleFunc("/getchat", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		result := make(map[string]interface{})
		result["chat_id"] = ""
		if err == nil {
			email, _ := admin.GetClaim(cookie.Value, "email")
			password, _ := admin.GetClaim(cookie.Value, "password")
			user := admin.GetUsers(client, database, "Users", bson.M{"email": email, "password": password}, bson.D{})[0]
			chat := chat.GetChats(client, database, "Chats", bson.M{"user_id": user.Id, "status": "Active"}, bson.D{})
			if len(chat) > 0 {
				result["chat_id"] = chat[0].ChatID
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/create_chat", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("auth_token")
		email, _ := admin.GetClaim(cookie.Value, "email")
		password, _ := admin.GetClaim(cookie.Value, "password")
		user := admin.GetUsers(client, database, "Users", bson.M{"email": email, "password": password}, bson.D{})[0]
		db := client.Database(database).Collection("Chats")
		chat := chat.Chat{UserID: user.Id, Status: "Active", Messages: []chat.Message{}, ChatID: time.Now().Unix()}
		db.InsertOne(r.Context(), chat)
		url := "http://127.0.0.1:8080/chat/" + strconv.FormatInt(chat.ChatID, 10)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		productsLimiter(w, r, client, database, collection, logs)
	})

	mux.HandleFunc("/products/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		products.Product(w, r, client, database, collection, logs)
	})
	cart.RegisterRoutes(mux, client, database, "Users")
	admin.RegisterRoutes(mux, client, database, "Users")
	auth.RegisterRoutes(mux, client, database, "Users")

	go chat.HandleMessages()

	handler := cors.Default().Handler(mux)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err = http.ListenAndServe(":8080", handler)
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
