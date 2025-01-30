package cart

import (
	"encoding/json"
	"net/http"
	"onlinestore/admin"
	"onlinestore/products"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"log"
	"html/template"
)

type CartPageData struct {
    Products []products.ProductModel
}

func AddCartHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodPost {
		product := products.ProductModel{}
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fmt.Println(product)
		token, _ := r.Cookie("auth_token")
		email, _ := admin.GetClaim(token.Value, "email")
		password, _ := admin.GetClaim(token.Value, "password")
		db := client.Database(database)
		err = addProductToUserProducts(context.Background(), db, email, password, product)
		if err != nil {
			log.Fatalf("Error adding product: %v", err)
		} else {
			log.Println("Product successfully added to user")
		}
	
		w.WriteHeader(http.StatusOK)
	}
}

func GetCartHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodGet {
		token, _ := r.Cookie("auth_token")
		email, _ := admin.GetClaim(token.Value, "email")
		password, _ := admin.GetClaim(token.Value, "password")
		user := admin.GetUsers(client, database, collection, bson.M{"email": email, "password": password}, bson.D{})
		tmpl, errTempl := template.ParseFiles("templates/cart.html")
		if errTempl != nil {
			http.Error(w, "Error with template", http.StatusInternalServerError)
		}
		data := CartPageData{
            Products: user[0].Products,
        }

		tmpl.Execute(w, data)
	}
}


func addProductToUserProducts(ctx context.Context, db *mongo.Database, email, password string, product products.ProductModel) error {
	collection := db.Collection("Users")

	filter := bson.M{"email": email, "password": password}
	update := bson.M{
		"$push": bson.M{
			"products": product,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add product to user's products: %v", err)
	}

	return nil
}