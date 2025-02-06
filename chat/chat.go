package chat

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"github.com/gorilla/websocket"
)

// Конфигурация WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Менеджер клиентов и сообщений
var clients = make(map[*websocket.Conn]bool) 
var broadcast = make(chan Message)

// Структура сообщения
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка подключения: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка сообщения: %v", err)
			}
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка отправки : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func GetChats(client *mongo.Client, database, collection string, filter bson.M, sorting bson.D) []Chat {
    coll := client.Database(database).Collection(collection)

    findOptions := options.Find().SetSort(sorting)

    cursor, err := coll.Find(context.TODO(), filter, findOptions)
    if err != nil {
        panic(err)
    }

    var chats []Chat
    if err := cursor.All(context.TODO(), &chats); err != nil {
        panic(err)
    }

    return chats
}