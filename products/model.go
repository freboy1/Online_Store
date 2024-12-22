package products

type ProductModel struct {
	ID          int64    `bson:"id"`
	Name        string `bson:"name"`
	Category    string `bson:"category"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	Discount    int    `bson:"discount"`
	Quantity    int    `bson:"quantity"`
}

type PageData struct {
	Products []ProductModel
	Error    string
}
