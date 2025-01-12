package products

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func TestInsertOne(t *testing.T) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to create Mongo client: %v", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer client.Disconnect(context.Background())

	database := "testdb"
	collection := "products"
	coll := client.Database(database).Collection(collection)

	product := ProductModel{
		Name:     "Product A",
		Price:    100,
		Quantity: 10,
	}

	_, err = insertOne(client, context.Background(), database, collection, product)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	var result ProductModel
	err = coll.FindOne(context.Background(), bson.M{"name": "Product A"}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find inserted product: %v", err)
	}

	assert.Equal(t, product.Name, result.Name, "Product name should match")
	assert.Equal(t, product.Price, result.Price, "Product price should match")
	assert.Equal(t, product.Quantity, result.Quantity, "Product quantity should match")
}
