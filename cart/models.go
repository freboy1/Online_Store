package cart

import (
	"onlinestore/products"
	"github.com/google/uuid"
)

type TransactionModel struct {
	transactionID int64    `bson:"id"`
	UserID        uuid.UUID    `bson:"user_id"`
	Products []products.ProductModel	`bson:"products"`
	TransactionStatus string `bson:"status"`
}