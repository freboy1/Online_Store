package products

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	if r.Method == http.MethodGet {
		GetProducts(client, database, collection)
	}
}


func GetProducts(client *mongo.Client, database, collection string) {
	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	var results []ProductModel
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	fmt.Println(results)
}
// func insertOne (client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

//     // select database and collection ith Client.Database method 
//     // and Database.Collection method
//     collection := client.Database(dataBase).Collection(col)
	
//     // InsertOne accept two argument of type Context 
//     // and of empty interface   
//     result, err := collection.InsertOne(ctx, doc)
//     return result, err
// }