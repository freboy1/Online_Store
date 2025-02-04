package cart

import (
	"github.com/gorilla/mux"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
)



func RegisterRoutes(router *mux.Router, client *mongo.Client, database, collection string) {
	router.HandleFunc("/addcart", func(w http.ResponseWriter, r *http.Request) {
		AddCartHandler(w, r, client, database, collection)
	})
	router.HandleFunc("/getcart", func(w http.ResponseWriter, r *http.Request) {
		GetCartHandler(w, r, client, database, collection)
	})
	router.HandleFunc("/verifycart", func(w http.ResponseWriter, r *http.Request) {
		VerifyCart(w, r, client, database, collection)
	})
}
