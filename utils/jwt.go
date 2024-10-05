package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt key used to create the signature for the JWT
var jwtSecret = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJwt generates a new JWT token
func GenerateJwt(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJwt parses a JWT token
func ParseJwt(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("Error parsing JWT token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Printf("Invalid JWT token")
		return nil, err
	}

	return claims, nil
}

// RefreshJwt refreshes a JWT token
func RefreshJwt(tokenString string) (string, error) {
	claims, err := ParseJwt(tokenString)
	if err != nil {
		return "", err
	}

	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ExtractClaims extracts the claims from a JWT token
func ExtractClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("Error extracting claims from JWT token: %v", err)
		return nil, err
	}

	return claims, nil
}

// ValidateJwt validates a JWT token
func ValidateJwt(tokenString string) bool {
	_, err := ParseJwt(tokenString)
	return err == nil
}

// Define a custom type for the context key
type contextKey string

const usernameKey contextKey = "username"

// SetUserInContext sets the username in the request context
func SetUserInContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

// GetUserFromContext gets the username from the request context
func GetUserFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

// Define the User type
type User struct {
	Username string `json:"username"`
	// Add other fields as needed
}

// SetUserToContext sets the user in the request context
func SetUserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, usernameKey, user)
}
