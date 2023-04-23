package account

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const OverTime = 48 // 小时
var jwtKey = []byte("my_secret_key")

type JwtCustClaims struct {
	Id   int
	Name string
	jwt.RegisteredClaims
}

func GenerateToken(id int, name string) string {
	MyJwtCustClaims := JwtCustClaims{
		Id:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(OverTime * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}
	tokenSetting := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJwtCustClaims)
	token, err := tokenSetting.SignedString(jwtKey)
	if err != nil {
		return err.Error()
	}
	return token
}
