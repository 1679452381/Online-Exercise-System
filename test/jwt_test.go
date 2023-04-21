package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

var jwtKey = []byte("zxczxc")

type UserClaim struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

// // 解析token
func AnalyToken(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Name, claims.RegisteredClaims.Issuer)
		return claims, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func TestJwt(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		Identity:         "1123123243245",
		Name:             "admin",
		RegisteredClaims: jwt.RegisteredClaims{},
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(tokenString)
	uc, err := AnalyToken(tokenString)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(uc)

}
