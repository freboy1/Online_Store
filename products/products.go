package products

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"fmt"
	"strconv"
	"html/template"
	"time"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	var pageData PageData
	products := GetProducts(client, database, collection, bson.M{}, bson.D{})
	pageData.Products = products
	if r.Method == http.MethodGet {
		r.ParseForm()
		filters := r.Form["filter"]
		if len(filters) != 0 {
			filter := bson.M{"category": bson.M{"$in": filters}}
			pageData.Products = GetProducts(client, database, collection, filter, bson.D{})
		}
		tmpl, err := template.ParseFiles("templates/products.html")
		if err != nil {
			pageData.Error = "Error with template"
		}
		tmpl.Execute(w, pageData)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			pageData.Error = fmt.Sprintf("ParseForm() error: %v", err)
			tmpl, _ := template.ParseFiles("templates/products.html")
			tmpl.Execute(w, pageData)
			return
		}
		name, description, priceStr, discountStr, quantityStr :=  r.FormValue("name"),  r.FormValue("description"),  r.FormValue("price"),  r.FormValue("discount"),  r.FormValue("quantity")
		product, err := checkProduct(name, description, priceStr, discountStr, quantityStr)
		if err != nil {
			pageData.Error = ("Input for price, discount, quantity must be numbers")
			tmpl, _ := template.ParseFiles("templates/products.html")
			tmpl.Execute(w, pageData)
			return
		}
		result, err := insertOne(client, context.TODO(), database, collection, product)
		if err != nil {
			pageData.Error = "Error inserting product: " + err.Error()
			tmpl, _ := template.ParseFiles("templates/products.html")
			tmpl.Execute(w, pageData)
			return

		}
	
		fmt.Println("Inserted product with ID:", result.InsertedID)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
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
func GetProducts(client *mongo.Client, database, collection string, filter bson.M, sorting bson.D) []ProductModel {
    // Access the collection from the specified database
    coll := client.Database(database).Collection(collection)

    // Define find options with sorting
    findOptions := options.Find().SetSort(sorting)

    // Fetch documents from the collection using the filter and find options
    cursor, err := coll.Find(context.TODO(), filter, findOptions)
    if err != nil {
        panic(err) // Handle error appropriately in production code
    }

    var products []ProductModel
    if err := cursor.All(context.TODO(), &products); err != nil {
        panic(err) // Handle error appropriately in production code
    }

    return products
}

func insertOne (client *mongo.Client, ctx context.Context, dataBase, col string, product ProductModel) (*mongo.InsertOneResult, error) {

    // select database and collection ith Client.Database method 
    // and Database.Collection method
    collection := client.Database(dataBase).Collection(col)
	product.ID = time.Now().Unix()
    // InsertOne accept two argument of type Context 
    // and of empty interface   
    result, err := collection.InsertOne(ctx, product)
    return result, err
}

func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col string, id int64) (*mongo.DeleteResult, error) {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"id", id}}
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}

func updateOne(client *mongo.Client, ctx context.Context, dataBase, col string, id int64, Product ProductModel) error {
	collection := client.Database(dataBase).Collection(col)
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"name", Product.Name},
			{"description", Product.Description},
			{"price", Product.Price},
			{"discount", Product.Discount},
			{"quantity", Product.Quantity},
		}},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("failed to update product")
		return err
	}

	// Check if the product was found and updated
	if result.MatchedCount == 0 {
		return err
	}

	fmt.Printf("Successfully updated %d product(s)\n", result.ModifiedCount)
	return nil
}