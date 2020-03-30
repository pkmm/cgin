package service

import (
	"cgin/conf"
	"cgin/model"
	"errors"
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
	Uid      uint64 `json:"uid"`
	UserName string `json:"user_name"`
	RoleIds  []int  `json:"role_ids"`
	jwt.StandardClaims
}

func (srv *jwtService) GetSignKey() []byte {
	return []byte(conf.AppConfig.String(conf.JwtSignKey))
}

func (srv *jwtService) GenerateToken(user *model.User) (string, error) {
	roleIds := make([]int, 0)
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.ID)
	}
	claims := AuthClaims{
		user.Id,
		user.OpenId,
		roleIds,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(conf.GetJwtExpiresAt()).Unix(),
			Issuer:    "c_gin",
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
	if token == nil {
		return nil, errors.New("token is not valid")
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
