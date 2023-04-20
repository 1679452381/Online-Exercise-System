package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// 生成md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

type UserClaim struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("zxczxc")

// 生成token
func GenerateToken(identity, name string) (string, error) {
	uc := &UserClaim{
		Identity:         identity,
		Name:             name,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// // 解析token
func AnalyToken(tokenString string) (*UserClaim, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
