package database

import (
	"database/sql"
	"fmt"
	"os"

	"monolith/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Initialize() error {
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_PORT"))
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return createTables()
}

func Close() {
	db.Close()
}

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT false
		);

		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price DECIMAL(10, 2) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS cart_items (
			user_id INTEGER REFERENCES users(id),
			product_id INTEGER REFERENCES products(id),
			quantity INTEGER NOT NULL,
			PRIMARY KEY (user_id, product_id)
		);
	`)
	return err
}

func CreateUser(user *models.User) error {
	return db.QueryRow("INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3) RETURNING id", 
		user.Username, user.Password, user.IsAdmin).Scan(&user.ID)
}

func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow("SELECT id, username, password, is_admin FROM users WHERE username = $1", 
		username).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}







func GetProductByID(id int) (*models.Product, error) {
	product := &models.Product{}
	err := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return product, nil
}


func SearchProducts(query string) ([]models.Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products WHERE LOWER(name) LIKE LOWER($1)", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func AddToCart(userID, productID, quantity int) error {
	_, err := db.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity", 
		userID, productID, quantity)
	return err
}

func Checkout(userID int) error {
	_, err := db.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)
	return err
}