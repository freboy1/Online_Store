package products

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/sirupsen/logrus"
	"onlinestore/logger"
)

func Product(w http.ResponseWriter, r *http.Request, client *mongo.Client, database, collection string, log *logrus.Logger) {
	vars := mux.Vars(r)
	idStr := vars["id"] 
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var pageData PageData

	product := GetProducts(client, database, collection, bson.M{"id": id}, bson.D{})
	pageData.Products = product
	pageData.Error = ""
	
	tmpl, err := template.ParseFiles("templates/product.html")
	if err != nil {
		pageData.Error = "Error with template"
	}

	if r.Method == http.MethodGet {
		logger.LogUserAction(log, product[0].Name, "1", idStr, map[string]interface{}{})
		tmpl.Execute(w, pageData)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			pageData.Error = fmt.Sprintf("ParseForm() error: %v", err)
			tmpl, _ := template.ParseFiles("templates/products.html")
			tmpl.Execute(w, pageData)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "delete":
			result, err := deleteOne(client, context.TODO(), database, collection, id)
			if err != nil {
				pageData.Error = "Could not delete"
				tmpl.Execute(w, pageData)
				return
			}
		
			fmt.Println("Deleted succesfully:", result)
			http.Redirect(w, r, "http://127.0.0.1:8080/products", http.StatusSeeOther)
			return
		case "update":
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
				tmpl, _ := template.ParseFiles("templates/product.html")
				tmpl.Execute(w, pageData)
				return
			}
			err = updateOne(client, context.TODO(), database, collection, id, product)
			if err != nil {
				pageData.Error = "Error inserting product: " + err.Error()
				tmpl, _ := template.ParseFiles("templates/product.html")
				tmpl.Execute(w, pageData)
				return
			}
			http.Redirect(w, r, "http://127.0.0.1:8080/products", http.StatusSeeOther)
			return
		}
	}
}