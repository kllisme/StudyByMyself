package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	jwtmiddleware "gopkg.in/iris-contrib/middleware.v5/jwt"
	iris "gopkg.in/kataras/iris.v5"
)

func SetupJWT() {

	isDevelopment := viper.GetBool("isDevelopment")
	secret := viper.GetString("server.jwt.secret")
	_jwt := jwtmiddleware.New(jwtmiddleware.Config{
		Debug: isDevelopment,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	JWT = _jwt

}

var (
	JWT           *jwtmiddleware.Middleware
	Authorization = func(ctx *iris.Context) {
		err := JWT.CheckJWT(ctx)
		if err != nil {
			Render(ctx, "000001", nil)
		} else {
			ctx.Next()
		}
	}
)
