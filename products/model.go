package products

import "time"

type ProductModel struct {
	ID          int64     `bson:"id"` // Unique identifier with an index
	Name        string    `bson:"name"` // Text index for search
	Category    string    `bson:"category"` // Category field for filtering and grouping
	Description string    `bson:"description"`
	Price       int       `bson:"price"` // Price field for sorting and filtering
	Discount    int       `bson:"discount"`
	Quantity    int       `bson:"quantity"`
	CreatedAt   time.Time `bson:"created_at"` // Timestamp for performance optimization (TTL index)
	Tags        []string  `bson:"tags"` // Array field for embedded queries
	Reviews     []Review  `bson:"reviews"` // Embedded array for customer feedback
}

type Review struct {
	UserID  int64  `bson:"user_id"`
	Rating  int    `bson:"rating"`
	Comment string `bson:"comment"`
	Date    time.Time `bson:"date"`
}

type PageData struct {
	Products []ProductModel
	Pages	[]int
	Error    string
}
