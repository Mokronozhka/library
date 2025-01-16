package server

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var key = []byte("test")

func CreateToken(UID string) (string, error) {

	payload := jwt.MapClaims{
		"iss": "Server",
		"sub": UID,
		"exp": time.Now().Add(5 * time.Minute),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, payload)

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func validateToken() {

}
