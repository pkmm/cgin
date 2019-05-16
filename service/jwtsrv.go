package service

import (
	"cgin/conf"
	"cgin/model"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

type jwtService struct {
	mutex *sync.Mutex
}

var JWTSrv = &jwtService{
	mutex: &sync.Mutex{},
}

type AuthClaims struct {
	Uid uint64 `json:"uid"`
	jwt.StandardClaims
}

func (srv *jwtService) GetSignKey() []byte {
	return []byte(conf.AppConfig.String("jwt_secret"))
}

func (srv *jwtService) GenerateToken(user *model.User) (string, error) {
	claims := AuthClaims{
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "ccla",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(srv.GetSignKey())
	return ss, err
}

func (srv *jwtService) GetAuthClaims(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.GetSignKey(), nil
	})
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
