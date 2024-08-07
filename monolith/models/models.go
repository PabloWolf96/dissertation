package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
    IsAdmin  bool   `json:"is_admin"`
}

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type CartItem struct {
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}

type Cart struct {
    UserID int        `json:"user_id"`
    Items  []CartItem `json:"items"`
}

type JWTClaim struct {
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
    jwt.StandardClaims
}