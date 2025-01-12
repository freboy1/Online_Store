package admin

import (
	"github.com/gorilla/mux"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
)

var (
	secretKey = []byte("secret -key")
)

func RegisterRoutes(router *mux.Router, client *mongo.Client, database, collection string) {
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if !adminMiddleware(w, r) {
			return
		}
		AdminPanelHandler(w, r)
	})
	router.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		if !adminMiddleware(w, r) {
			return
		}
		UsersHandler(w, r, client, database, collection)
	})
	router.HandleFunc("/admin/user/{id:[0-9a-fA-F-]+}", func(w http.ResponseWriter, r *http.Request) {
		if !adminMiddleware(w, r) {
			return
		}
		UserHandler(w, r, client, database, collection)
	})
	router.HandleFunc("/admin/send-email", func(w http.ResponseWriter, r *http.Request) {
		if !adminMiddleware(w, r) {
			return
		}
		SendEmailHandler(w, r, client, database, collection)
	})
}

func adminMiddleware(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("auth_token")
	isadmin := false
	if err == nil {
		role, _ := GetRole(cookie.Value)
		if role == "admin" {
			isadmin = true 
		}
	}
	if !isadmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func GetRole(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	value, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found in token")
	}

	return value, nil

}