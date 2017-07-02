package service

import (
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"gopkg.in/kataras/iris.v5"
)

type TokenService struct {
}

func (self *TokenService) Token(ctx *iris.Context) (string, error) {
	secret := viper.GetString("server.jwt.secret")
	tokenExpire := time.Duration(viper.GetInt("server.jwt.tokenExpire"))
	// cookieExpire := time.Duration(viper.GetInt("server.jwt.cookieExpire"))
	cookieName := viper.GetString("server.jwt.cookieName")
	// cookieDomain := viper.GetString("server.jwt.cookieDomain")
	issuer := viper.GetString("server.jwt.issuer")
	_tokenExpire := time.Now().Add(time.Hour * tokenExpire).Unix()
	// _cookieExpire := time.Now().Add(time.Hour * cookieExpire)
	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionId": ctx.Session().ID(),
		"exp":       _tokenExpire,
		"iss":       issuer,
	})
	token, err := _token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	ctx.SetCookieKV(cookieName, token)
	return token, nil
}

func (self *TokenService) GetToken(ctx *iris.Context) string {
	token := ""
	auth := ctx.RequestHeader("Authorization")
	if auth != "" {
		auth = strings.TrimLeft(auth, " ")
		auth = strings.TrimRight(auth, " ")
		ss := strings.Split(auth, " ")
		if len(ss) == 2 {
			token = ss[1]
		}
	}
	return token
}
