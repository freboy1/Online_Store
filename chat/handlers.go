package chat

import (
	"encoding/json"
	"context"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
)

// Хранилище истории сообщений



func HandleMessageHistory(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Println("Request body:", string(body))

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}
	
	intID, _ := strconv.ParseInt(data["id"].(string), 10, 64)
	if r.Method == "POST" {
		var message Message
		
		message.Username, _ = data["username"].(string)
		message.Content, _ = data["content"].(string)

		coll := client.Database(database).Collection(collection)
		filter := bson.M{"id": intID}
		update := bson.M{
			"$push": bson.M{"messages": message},
		}

		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "Failed to update chat", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}


func HandleGetMessageHistory(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string) {
	vars := mux.Vars(r)
	idStr := vars["id"] 
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println(id)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	messageHistory := GetChats(client, database, collection, bson.M{"id": id, "status": "Active"}, bson.D{})[0].Messages
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageHistory)
}