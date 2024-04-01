package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("joseph-chatroom")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, false
	}
	if !token.Valid {
		return nil, false
	}
	return claims, true
}

var (
	ErrNoAuthToken            = errors.New("no authentication token provided")
	ErrInvalidAuthTokenFormat = errors.New("invalid authentication token format")
	ErrInvalidToken           = errors.New("invalid or expired token")
)

func AuthenticateUser(r *http.Request) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrNoAuthToken
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return nil, ErrInvalidAuthTokenFormat
	}

	claims, valid := ValidateToken(bearerToken[1])
	if !valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
