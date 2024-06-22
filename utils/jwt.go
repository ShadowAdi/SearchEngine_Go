package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaim struct {
	Id    string `json:"id"`
	User  string `json:"user"`
	Admin bool   `json:"role"`
	jwt.RegisteredClaims
}

func CreateNewAuthToken(id string, email string, isAdmin bool) (string, error) {
	claims := AuthClaim{
		Id:    id,
		User:  email,
		Admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "searchengine.com",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret_key, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		panic("Can't Find Secret Key")
	}
	signedToken, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", errors.New("Error signing the token")
	}
	return signedToken, nil

}
