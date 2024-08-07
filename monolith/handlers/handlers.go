package handlers

import (
	"encoding/json"
	"monolith/database"
	"monolith/models"
	"monolith/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	err := database.CreateUser(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, user)
}

func Login(w http.ResponseWriter, r *http.Request) {
    var loginUser models.User
    json.NewDecoder(r.Body).Decode(&loginUser)

    user, err := database.GetUserByUsername(loginUser.Username)
    if err != nil {
        utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    if !utils.CheckPasswordHash(loginUser.Password, user.Password) {
        utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    token, err := utils.GenerateJWT(user)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}




func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	product, err := database.GetProductByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	products, err := database.SearchProducts(query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(products) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No products found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
    var cartItem models.CartItem
    json.NewDecoder(r.Body).Decode(&cartItem)

    // Get the user ID from the context
    userID, ok := r.Context().Value("userID").(int)
    if !ok {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user ID")
        return
    }

    err := database.AddToCart(userID, cartItem.ProductID, cartItem.Quantity)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Item added to cart"})
}

func Checkout(w http.ResponseWriter, r *http.Request) {
    // Get the user ID from the context
    userID, ok := r.Context().Value("userID").(int)
    if !ok {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user ID")
        return
    }

    err := database.Checkout(userID)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Checkout successful"})
}
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "server is healthy"})
}