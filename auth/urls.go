package auth

import (
	"github.com/gorilla/mux"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *mux.Router, client *mongo.Client, database, collection string) {
	router.HandleFunc("/login", LoginGet).Methods("GET")
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        LoginHandler(w, r, client, database, collection)
    })
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        RegisterHandler(w, r, client, database, collection)
    })
	router.HandleFunc("/register/email-verification", func(w http.ResponseWriter, r *http.Request) {
        EmailVerification(w, r, client, database, collection)
    })
}
