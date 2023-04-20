package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

var jwtKey = []byte("zxczxc")

func TestJwt(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(tokenString)

}
