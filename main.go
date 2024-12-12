package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"onlinestore/db"
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

type ProductModel struct {
	Name string
	Description string
	Price int
	Discount int
	Quantity int
}


func getProducts(w http.ResponseWriter, r *http.Request) {

}


func main() {
	uri := "mongodb://localhost:2707"
	client, ctx, dbErr := db.Connect(uri)
	if dbErr {
		log.Println("Could not connect to the MongoDB")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/products", getProducts)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}