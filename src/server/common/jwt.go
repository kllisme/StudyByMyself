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

		//token首先从Cookie里获取，若是没有信息，再尝试从请求头里获取
		// Extractor: func(ctx *iris.Context) (string, error) {

		// 	tokenString := ctx.GetCookie("Authorization")
		// 	if tokenString == "" {
		// 		return jwtmiddleware.FromAuthHeader(ctx)
		// 	}
		// 	return tokenString, nil
		// },
	})
	JWT = _jwt

}

var (
	JWT           *jwtmiddleware.Middleware
	Authorization = func(ctx *iris.Context) {
		err := JWT.CheckJWT(ctx)
		if err != nil {
			Render(ctx, "000001", err)
		} else {
			ctx.Next()
		}
	}
)
