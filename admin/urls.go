package admin

import (
	"github.com/gorilla/mux"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *mux.Router, client *mongo.Client, database, collection string) {
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		AdminPanelHandler(w, r)
	})
	router.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		UsersHandler(w, r, client, database, collection)
	})
	router.HandleFunc("/admin/send-email", func(w http.ResponseWriter, r *http.Request) {
		SendEmailHandler(w, r, client, database, collection)
	})
}
