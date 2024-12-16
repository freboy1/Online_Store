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




func main() {
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
        products.ProductsHandler(w, r, client, database, collection)
    })
	mux.HandleFunc("/products/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
        products.Product(w, r, client, database, collection)
    })

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}