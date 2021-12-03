package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = "s3c4e7" //ofc this should be on the config

func CreateJWT() (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		Issuer:    "test",
		Subject:   "auth0|test@example.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

//CheckJWT base on https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/
func CheckJWT(stringToken string) error {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if x, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Sprintln(x)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.Parse(stringToken, keyFunc)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
