package products

import (
	"net/http"
	"net/http/httptest"
	"onlinestore/db"
	"testing"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"strings"
)

func TestGetProduct(t *testing.T) {

	router := mux.NewRouter()
	client, err := db.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	router.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		Product(w, r, client, "onlineStore", "AlisherExpress", nil)
	})

	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}

	name, description, price, discount, quantity := extractValues(response.Body.String())
	
	if name != "Smartphone" || description != "High-end smartphone with 128GB storage." || price != "699" || discount != "10" || quantity != "50" {
		t.Errorf("Incorrect response body, Got: %s", response.Body.String())
	}
}


func extractValues(html string) (string, string, string, string, string) {
	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", "", "", ""
	}

	// Extract values
	name := doc.Find("input#name").AttrOr("value", "")
	description := doc.Find("input#description").AttrOr("value", "")
	price := doc.Find("input#price").AttrOr("value", "")
	discount := doc.Find("input#discount").AttrOr("value", "")
	quantity := doc.Find("input#quantity").AttrOr("value", "")
	return name, description, price, discount, quantity
}