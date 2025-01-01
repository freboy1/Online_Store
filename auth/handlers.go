package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("The request body is %v\n", r.Body)

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request body")
		return
	}

	fmt.Printf("The user request value %v\n", u)

	if u.Username == "Chek" && u.Password == "123456" {
		tokenString, err := CreateToken(u.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error creating token")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Login successful")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "Invalid credentials")
}


func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized: Missing or invalid token")
		return
	}

	err = verifyToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized: Invalid token")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to the protected area")
}
