package utils

import (
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
func GenerateJwt(username, role string) (string, error) {
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
