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
		fmt.Println(err)
		return "null"
	}
	return token
}

func ChecKToken(token string) (bool, int) {
	for _, user := range Users {
		if user.Token == token {
			id, accessPass := ParseJWTToken(token)
			if id == user.Id && accessPass {
				return true, id
			}

		}
	}

	return false, -1
}
