package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var key = []byte("SecretKey")

func CreateToken(UID string) (string, error) {

	//payload := jwt.MapClaims{
	//	"iss": "Server",
	//	"sub": UID,
	//	"exp": time.Now().Add(5 * time.Minute),
	//}

	payload := jwt.RegisteredClaims{
		Issuer:    "Server",
		Subject:   UID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateToken(tokenString string) error {

	//var payload jwt.MapClaims

	//token, err := jwt.ParseWithClaims(tokenString, &payload, func(token *jwt.Token) (interface{}, error) {
	//	return key, nil
	//})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return err

	}

	if !token.Valid {
		return errors.New("token is invalid")
	}

	return nil

}
