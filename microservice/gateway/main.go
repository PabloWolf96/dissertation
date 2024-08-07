package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/signup", forwardToAuth).Methods("POST")
	r.HandleFunc("/login", forwardToAuth).Methods("POST")
	r.HandleFunc("/products/search", forwardToProduct).Methods("GET")
	r.HandleFunc("/product/{id}", forwardToProduct).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Protected routes
	r.HandleFunc("/cart", forwardToCart).Methods("POST")
	r.HandleFunc("/cart/checkout", forwardToCart).Methods("POST")

	log.Println("Gateway is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func forwardToAuth(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, os.Getenv("AUTH_SERVICE_URL"))
}

func forwardToProduct(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, os.Getenv("PRODUCT_SERVICE_URL"))
}

func forwardToCart(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, os.Getenv("CART_SERVICE_URL"))
}

func forwardRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
    log.Printf("Forwarding request to %s%s", serviceURL, r.URL.Path)
    client := &http.Client{}
    req, err := http.NewRequest(r.Method, serviceURL+r.URL.Path+"?"+r.URL.RawQuery, r.Body)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    req.Header = r.Header
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    log.Printf("Received response with status code %d", resp.StatusCode)

    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        http.Error(w, "Error reading response body", http.StatusInternalServerError)
        return
    }
    log.Printf("Response body: %s", string(body))
    _, err = w.Write(body)
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gateway is healthy"))
}