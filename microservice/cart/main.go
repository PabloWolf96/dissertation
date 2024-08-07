package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type JWTClaim struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    IsAdmin  bool   `json:"is_admin"`
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/cart", addToCart).Methods("POST")
	r.HandleFunc("/cart/checkout", checkout).Methods("POST")

	log.Println("Cart service is running on port 8003")
	log.Fatal(http.ListenAndServe(":8003", r))
}

func initDB() {
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_PORT"))
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cart_items (
			user_id INTEGER,
			product_id INTEGER,
			quantity INTEGER NOT NULL,
			PRIMARY KEY (user_id, product_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claims, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	var cartItem CartItem
	json.NewDecoder(r.Body).Decode(&cartItem)

	_, err = db.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity",
    claims.UserID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
}

func checkout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	claims, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("DELETE FROM cart_items WHERE user_id = $1", claims.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Checkout successful"})
}

func validateToken(token string) (*JWTClaim, error) {
    authServiceURL := os.Getenv("AUTH_SERVICE_URL")
    
    req, err := http.NewRequest("POST", authServiceURL+"/validate", nil)
    if err != nil {
        return nil, err
    }
    
    // Forward the Authorization header as-is
    req.Header.Set("Authorization", token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to validate token: status %d", resp.StatusCode)
    }

    var claims JWTClaim
    if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
        return nil, err
    }

    return &claims, nil
}