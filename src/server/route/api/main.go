package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/route/api/admin"
)

func Api(app *iris.Framework) {

	var (
		user    = &api.UserController{}
		public  = &api.PublicController{}
		captcha = &api.CaptchaController{}
		login   = &api.LoginController{}
	)
	v1 := app.Party("/v1", func(ctx *iris.Context) {
		ctx.Next()
	})
	v1.Get("/captcha.png", captcha.Captcha)

	v1.UseFunc(common.CORS.Serve)
	v1.Options("/login",common.CORS.Serve)

	v1.Get("/token", public.Token)
	v1.Post("/login", login.Login)

	v1.UseFunc(common.Authorization)

	admin.Setup(v1)

	userApi := v1.Party("/user")
	{
		userApi.Get("/detail", user.GetSessionUser)
	}

}
