package middleware

import (
	"context"
	"net/http"
	"strings"

	"monolith/database"
	"monolith/utils"
)

func JwtAuthentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            utils.RespondWithError(w, http.StatusUnauthorized, "Missing auth token")
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
            return
        }

        claims, err := utils.ValidateToken(bearerToken[1])
        if err != nil {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
            return
        }

        user, err := database.GetUserByUsername(claims.Username)
        if err != nil {
            utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
            return
        }

        // Add the user ID to the request context
        ctx := context.WithValue(r.Context(), "userID", user.ID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}