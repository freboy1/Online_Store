package products

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Product(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	vars := mux.Vars(r)
	idStr := vars["id"] 
	id, _ := strconv.ParseInt(idStr, 10, 64)
	fmt.Println(id)
	if r.Method == http.MethodGet {
		product := GetProducts(client, database, collection, bson.D{{"id", id}})
		fmt.Println(product)
	}
}