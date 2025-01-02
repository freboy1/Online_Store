package models

import (
	"onlinestore/products"
)

type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	Code string `json:"code"`
	Verified string `json:"verified"`
	Cash int `json:"cash"`
	Products []products.ProductModel `json:"products"`
}

type Message struct {
	Subject string
	Text string
}