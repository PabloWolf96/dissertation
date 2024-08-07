package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}

type JWTClaim struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/signup", signUp).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/validate", validateToken).Methods("POST")

	log.Println("Auth service is running on port 8001")
	log.Fatal(http.ListenAndServe(":8001", r))
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
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT false
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	err := db.QueryRow("INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.IsAdmin).Scan(&user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginUser User
	json.NewDecoder(r.Body).Decode(&loginUser)

	var user User
	err := db.QueryRow("SELECT id, username, password, is_admin FROM users WHERE username = $1",
		loginUser.Username).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !checkPasswordHash(loginUser.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(&user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func validateToken(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Missing auth token", http.StatusUnauthorized)
        return
    }

    bearerToken := strings.Split(authHeader, " ")
    if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
        http.Error(w, "Invalid token format", http.StatusUnauthorized)
        return
    }

    tokenString := bearerToken[1]

    claims, err := validateJWT(tokenString)
    if err != nil {
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(claims)
}

func generateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Username: user.Username,
		UserID:   user.ID,  
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateJWT(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, err
	}

	return claims, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}