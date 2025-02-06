package chat

import (
	"github.com/google/uuid"
)

type Chat struct {
	UserID uuid.UUID `bson:"user_id"` 
	ChatID int64    `bson:"id"`
	Messages []Message `bson:"messages"`
	Status string `bson:"status"`
}