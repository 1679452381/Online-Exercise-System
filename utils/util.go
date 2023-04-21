package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
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
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Name, claims.RegisteredClaims.Issuer)
		return claims, nil
	} else {
		return nil, err
	}
}

// 获取uuid
func GetUUID() string {
	return fmt.Sprintf("%x", uuid.New())
}
