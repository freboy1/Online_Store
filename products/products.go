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
	"github.com/sirupsen/logrus"
	"onlinestore/logger"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string, log *logrus.Logger) {
	var pageData PageData
	products := GetProducts(client, database, collection, bson.M{}, bson.D{})
	pageData.Products = products
	if r.Method == http.MethodGet {
		r.ParseForm()
		page, err := getPage(r.URL.Query().Get("page"))
		if err != nil || page == 0 {
			pageData.Error = "Pagination must be numbers"
		}
		products, err = filterSortProducts(products, client, database, collection, r.Form["filter"], r.FormValue("sort"))
		if err != nil {
			pageData.Error = "Mistake with sort"
		}
		pageData.Pages, pageData.Products, err = Paginate(products, 3, page)
		if err != nil {
			pageData.Error = "Wrong number of pagination"
		}
		tmpl, errTempl := template.ParseFiles("templates/products.html")
		if errTempl != nil {
			pageData.Error = "Error with template"
		}
		logger.LogUserAction(log, "get Products", "1", "produtcs", map[string]interface{}{"filter": r.Form["filter"], "sort": r.FormValue("sort"), "Page Pagination": page})
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

func getPage(pageStr string) (int, error) {
	var page int
	if pageStr != "" {
		pageValue, err := strconv.Atoi(pageStr)
		if err != nil {
			return page, fmt.Errorf("Pagination must be numbers")
		} else {
			page = pageValue
		}
	} else {
		page = 1
	}
	return page, nil
}



func Paginate(products []ProductModel, step, page int) ([]int, []ProductModel, error) {
	lengthProducts := len(products)
	pages := make([]int, 0)
	if lengthProducts % 3 == 0 {
		pages = GeneratePages(lengthProducts / 3)
	} else {
		pages = GeneratePages((lengthProducts / 3) + 1)
	}

	start := (page * 3) - 3
	end := page * 3

	if start < 0 || start >= len(products) {
		return pages, products, fmt.Errorf("Invalid pagination")
	}
	if end > len(products) {
		end = len(products)
	}

	return pages, products[start:end], nil
}


func filterSortProducts(products []ProductModel, client *mongo.Client, database, collection string, filters []string, sortStr string) ([]ProductModel, error) {
	sort, err := strconv.Atoi(sortStr)
	if sortStr != "" {
		if err != nil {
			return products, err
		}
	}
	sorting := bson.D{{"price", sort}}
	if sort == 0 {
		sorting = bson.D{}
	}

	if len(filters) != 0 {
		filter := bson.M{"category": bson.M{"$in": filters}}
		products = GetProducts(client, database, collection, filter, sorting)
	} else if sort != 0 {
		products = GetProducts(client, database, collection, bson.M{}, sorting)
	}
	return products, nil
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

    collection := client.Database(dataBase).Collection(col)
	product.ID = time.Now().Unix()
 
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

func GeneratePages(n int) []int {
    pages := make([]int, n)
    for i := 0; i < n; i++ {
        pages[i] = i + 1
    }
    return pages
}