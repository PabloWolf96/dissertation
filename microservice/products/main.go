package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/product/{id}", getProduct).Methods("GET")
	r.HandleFunc("/products/search", searchProducts).Methods("GET")

	log.Println("Product service is running on port 8002")
	log.Fatal(http.ListenAndServe(":8002", r))
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
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price DECIMAL(10, 2) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var product Product
	err := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    log.Printf("Received search query: %s", query)
    if query == "" {
        http.Error(w, "Search query is required", http.StatusBadRequest)
        return
    }

    rows, err := db.Query("SELECT id, name, price FROM products WHERE LOWER(name) LIKE LOWER($1)", "%"+query+"%")
    if err != nil {
        log.Printf("Database query error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var product Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
            log.Printf("Error scanning row: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        products = append(products, product)
    }

    log.Printf("Found %d products", len(products))

    if len(products) == 0 {
        http.Error(w, "No products found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(products)
}