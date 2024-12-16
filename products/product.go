package products

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"fmt"
	"context"
	"log"
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
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "delete":
			result, err := deleteOne(client, context.TODO(), database, collection, id)
			if err != nil {
				log.Fatal(err)
			}
		
			fmt.Println("Deleted succesfully:", result)
			http.Redirect(w, r, "http://127.0.0.1:8080/products", http.StatusSeeOther)
			return
		}
	}
}