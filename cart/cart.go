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

type Response struct {
	Email           string `bson:"email"`
    Password        string `bson:"password"`
	Status 			string	`bson:"status"`
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
		transaction_id, _ := addTransaction(client, context.Background(), database, "Transactions", user[0].Id, user[0].Products)
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
				"transactionId": transaction_id, 
			},
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка кодирования JSON:", err)
			return
		}

		url := "http://127.0.0.1:8081/buy"
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

		http.Redirect(w, r, "http://127.0.0.1:8081/payment", http.StatusSeeOther)
	}
}

type TransactionPage struct {
	TransactionStatus string `bson:"status"`
	Products []products.ProductModel	`bson:"products"`
}

func VerifyCart(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, col string) {
	if r.Method == http.MethodGet {
		token, _ := r.Cookie("auth_token")
		email, err := admin.GetClaim(token.Value, "email")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		password, _ := admin.GetClaim(token.Value, "password")
		user_id := admin.GetUsers(client, database, col, bson.M{"email": email, "password": password}, bson.D{})[0].Id
		col = "Transactions"
		dataBase := client.Database(database).Collection(col)

		result := dataBase.FindOne(context.Background(), bson.M{"user_id": user_id})
		var transactionPage TransactionPage
		_ = result.Decode(&transactionPage)

		tmpl, errTempl := template.ParseFiles("templates/cart-status.html")
		if errTempl != nil {
			http.Error(w, "Error with template", http.StatusInternalServerError)
		}
		data := TransactionPage{
            Products: transactionPage.Products,
			TransactionStatus: transactionPage.TransactionStatus,
        }
		
		tmpl.Execute(w, data)

	} else if r.Method == http.MethodPost {
		var response Response
		if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		collection := client.Database(database).Collection("Transactions")
		user_id := admin.GetUsers(client, database, col, bson.M{"email": response.Email, "password": response.Password}, bson.D{})[0].Id
		_, err := collection.UpdateOne(context.Background(), bson.M{"user_id": user_id}, bson.M{"$set": bson.M{"status": response.Status}})
		if err != nil {
			http.Error(w, "Failed to update document", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
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

func addTransaction (client *mongo.Client, ctx context.Context, dataBase, col string, user_id uuid.UUID, products []products.ProductModel) (int64, error) {
    collection := client.Database(dataBase).Collection(col)
	transaction := TransactionModel{
        TransactionStatus: "pending",
        Products:          products,
        UserID:            user_id,
        TransactionID:     time.Now().Unix(), // Proper assignment
    }
    _, err := collection.InsertOne(ctx, transaction)
    return transaction.TransactionID, err
}