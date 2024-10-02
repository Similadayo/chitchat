package middlewares

import (
	"fmt"
	"github/similadayo/chitchat/utils"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is empty
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		tokenStrings := strings.Split(authHeader, " ")[1]
		if tokenStrings == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Parse the JWT token
		claims, err := utils.ParseJwt(tokenStrings)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the username in the request context
		r = r.WithContext(utils.SetUserInContext(r.Context(), claims.Username))

		// Proceed to the next middleware or handler
		fmt.Println("User is authenticated", claims.Username)
		next.ServeHTTP(w, r)

	})
}

// CorsMiddleware is a middleware that adds CORS headers to the response
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ContentTypeMiddleware is a middleware that sets the Content-Type header to application/json
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
