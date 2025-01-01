package auth

import (
	"onlinestore/products"
)

type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Verified string `json:"verified"`
	Cash int `json:"cash"`
	Products []products.ProductModel `json:"products"`
}