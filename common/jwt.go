package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

const secret = "1P&dG^5MceRb0T#7QDu6OtF%)$Nh@q"

func VerifyToken(tokenStr string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}
