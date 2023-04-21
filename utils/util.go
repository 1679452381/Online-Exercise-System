package utils

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
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

// 校验邮箱格式
func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 生成4位数验证码
func GetCode() string {
	//初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	var code string
	for i := 0; i < 4; i++ {
		//strconv.Itoa  将整数型数据转换为字符串型数据
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// 发送邮箱验证码
func SendEmailCode(toEmail, code string) error {
	e := email.NewEmail()
	e.From = "GET <17700611471@163.com>"
	e.To = []string{toEmail}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您在注册在线练习系统，您的验证码为：<b>" + code + "</b>")
	//err := e.Send("smtp.163.com:465", smtp.PlainAuth("", "15660589213@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"))
	//返回EOF 关闭SSL重试
	return e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "17700611471@163.com", "DSBZHQSKFWQVDSVK", "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})

}
