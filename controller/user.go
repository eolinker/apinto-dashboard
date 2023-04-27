package controller

import (
	"errors"

	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	sessionCache user.ISessionCache
)

func init() {
	bean.Autowired(&sessionCache)
}

var jwtSecret = []byte("apintp-dashboard")

func SetUser(ginCtx *gin.Context) {
	cookie, err := ginCtx.Cookie(Session)
	if err != nil {
		return
	}

	session, _ := sessionCache.Get(ginCtx, sessionCache.Key(cookie))
	if session == nil {
		return
	}
	uc := UserClaim{}
	parse, err := jwt.ParseWithClaims(session.Jwt, &uc, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return
	}
	if parse.Valid {
		ginCtx.Set(UserId, uc.Id)
	}
}

func JWTEncode(claim *UserClaim) (string, error) {

	userClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	userJWT, err := userClaim.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return userJWT, nil
}

func JWTDecode(info string) (*UserClaim, error) {
	uc := UserClaim{}
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
