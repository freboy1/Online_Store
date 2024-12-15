package products

type ProductModel struct {
	ID          int64    `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	Discount    int    `bson:"discount"`
	Quantity    int    `bson:"quantity"`
}
