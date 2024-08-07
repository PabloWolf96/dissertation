package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"monolith/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)



var jwtKey = []byte("your_secret_key")

func GenerateJWT(user *models.User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &models.JWTClaim{
        Username: user.Username,
        IsAdmin:  user.IsAdmin,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(signedToken string) (*models.JWTClaim, error) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &models.JWTClaim{},
        func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        },
    )

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*models.JWTClaim)
    if !ok {
        return nil, errors.New("couldn't parse claims")
    }

    if claims.ExpiresAt < time.Now().Local().Unix() {
        return nil, errors.New("token expired")
    }

    return claims, nil
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}