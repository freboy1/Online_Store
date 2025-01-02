package auth

import (
	"github.com/gorilla/mux"
	// "net/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *mux.Router, client *mongo.Client, database, collection string) {
	router.HandleFunc("/login", LoginGet).Methods("GET")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
}
