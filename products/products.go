package products

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"fmt"
	"strconv"
	"log"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodGet {
		products := GetProducts(client, database, collection)
		w.Write(products)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name, description, priceStr, discountStr, quantityStr :=  r.FormValue("name"),  r.FormValue("desc"),  r.FormValue("price"),  r.FormValue("discount"),  r.FormValue("quantity")
		product, err := checkProduct(name, description, priceStr, discountStr, quantityStr)
		if err != nil {
			fmt.Println("Error with product")
			return
		}
		result, err := insertOne(client, context.TODO(), database, collection, product)
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Println("Inserted product with ID:", result.InsertedID)
	} else if r.Method == http.MethodDelete {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		result, err := deleteOne(client, context.TODO(), database, collection, name)
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Println("Deleted succesfully:", result)
	}
}

func checkProduct(name, description, priceStr, discountStr, quantityStr string) (ProductModel, error) {
	product := ProductModel{
		Name: name,
		Description: description,
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return product, err
	}
	discount, err := strconv.Atoi(discountStr)
	if err != nil {
		return product, err
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return product, err
	}
	product.Price = price
	product.Discount = discount
	product.Quantity = quantity
	return product, nil
}
func GetProducts(client *mongo.Client, database, collection string) []byte {
	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	var bsonProducts []ProductModel
	if err = cursor.All(context.TODO(), &bsonProducts); err != nil {
		panic(err)
	}
	jsonProducts, err := json.Marshal(bsonProducts)
	return jsonProducts
}
func insertOne (client *mongo.Client, ctx context.Context, dataBase, col string, product ProductModel) (*mongo.InsertOneResult, error) {

    // select database and collection ith Client.Database method 
    // and Database.Collection method
    collection := client.Database(dataBase).Collection(col)
	
    // InsertOne accept two argument of type Context 
    // and of empty interface   
    result, err := collection.InsertOne(ctx, product)
    return result, err
}

func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col, name string) (*mongo.DeleteResult, error) {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"name", name}}
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}