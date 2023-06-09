package account

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const OverTime = 24 // 小时
var jwtKey = []byte("my_secret_key")

type JwtCustClaims struct {
	id   int
	name string
	jwt.RegisteredClaims
}

// ParseJWTToken 解析 Token,获得过期时间
func ParseJWTToken(tokenString string) (int, bool) {
	// Parse the token string
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil // Replace with your own secret key
	})

	// Get the claims from the token
	claims, _ := token.Claims.(jwt.MapClaims)

	exp, _ := claims["exp"].(float64)
	idString, _ := claims["sub"].(string)
	id, _ := strconv.Atoi(idString)

	if exp < float64(time.Now().Unix()) {
		return -1, false
	}

	return id, true
}

/*
	Token 包含生成时间、过期时间、用户 id
*/

// GenerateToken 生成 Token
func GenerateToken(id int, name string) string {
	MyJwtCustClaims := JwtCustClaims{
		id:   id,
		name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(OverTime * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(id),
		},
	}
	tokenSetting := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJwtCustClaims)
	token, err := tokenSetting.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err) // ---日志
		return "null"
	}
	return token
}

// ChecKToken 检测 Token 是否过期
func ChecKToken(token string) (bool, int) {
	for _, user := range Users {
		if user.Token == token {
			id, accessPass := ParseJWTToken(token)
			if id == user.Id && accessPass {
				return true, id
			}

		} else {
			continue
		}
	}

	return false, -1
}
