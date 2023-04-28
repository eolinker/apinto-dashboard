package common

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eolinker/apinto-dashboard/controller"
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

func JWTEncode(claim *controller.UserClaim, jwtSecret interface{}) (string, error) {

	userClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	userJWT, err := userClaim.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return userJWT, nil
}

func JWTDecode(info string, jwtSecret interface{}) (*controller.UserClaim, error) {
	uc := controller.UserClaim{}
	parse, err := jwt.ParseWithClaims(info, &uc, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if parse.Valid {
		return &uc, nil
	}
	return nil, errors.New("invalid user info")
}
