package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"onlinestore/admin"
	"onlinestore/products"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	token, _ := r.Cookie("auth_token")
	email, _ := admin.GetClaim(token.Value, "email")
	password, _ := admin.GetClaim(token.Value, "password")
	user := admin.GetUsers(client, database, collection, bson.M{"email": email, "password": password}, bson.D{})
	if r.Method == http.MethodGet {
		tmpl, errTempl := template.ParseFiles("templates/cart.html")
		if errTempl != nil {
			http.Error(w, "Error with template", http.StatusInternalServerError)
		}
		data := CartPageData{
            Products: user[0].Products,
        }

		tmpl.Execute(w, data)
	} else if r.Method == http.MethodPost {
		addTransaction(client, context.Background(), database, "Transactions", user[0].Id, user[0].Products)
		cartItems := []map[string]interface{}{}
		for _, product := range user[0].Products {
			cartItem := map[string]interface{}{"id": product.ID, "name": product.Name, "price": product.Price}
			cartItems = append(cartItems, cartItem)
		}
		data := map[string]interface{}{
			"cartItems": cartItems,
			"customer": map[string]interface{}{
				"id": user[0].Id,
				"name": user[0].Username,
				"email": user[0].Email,
			},
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка кодирования JSON:", err)
			return
		}

		url := "http://localhost:8081/buy"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка создания запроса:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			return
		}
		defer resp.Body.Close()

		http.Redirect(w, r, "http://localhost:8081/payment", http.StatusSeeOther)
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

func addTransaction (client *mongo.Client, ctx context.Context, dataBase, col string, user_id uuid.UUID, products []products.ProductModel) (*mongo.InsertOneResult, error) {
    collection := client.Database(dataBase).Collection(col)
	transaction := TransactionModel{TransactionStatus: "pending", Products: products, UserID: user_id}
    transaction.transactionID = time.Now().Unix()
    result, err := collection.InsertOne(ctx, transaction)
    return result, err
}