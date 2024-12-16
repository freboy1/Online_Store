package products

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
)

func Product(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	vars := mux.Vars(r)
	idStr := vars["id"] 
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		product := GetProducts(client, database, collection, bson.D{{"id", id}})
		tmpl, err := template.ParseFiles("templates/product.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, product)
	}
}