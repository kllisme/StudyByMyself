package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/middleware"
)

func Api(app *iris.Framework) {

	var (
		user = &api.UserController{}
		//public = &api.PublicController{}
		captcha = &api.CaptchaController{}
		login = &api.LoginController{}
	)
	v1 := app.Party("/v1", func(ctx *iris.Context) {
		ctx.Next()
	})
	{
		v1.Get("/captcha.png", captcha.Captcha)

		//为跨域请求设定入口
		v1.UseFunc(common.CORS.Serve)
		v1.Options("/*anything", common.CORS.Serve)

		//v1.Get("/token", public.Token)
		v1.Post("/login", login.Login)

		v1.UseFunc(common.Authorization)
		v1.Get("/session/info", user.GetSessionInfo)

		//进行权限控制的接口
		AccessControlledAPI := v1.UseFunc(middleware.AccessControlMiddleware)
		{
			AccessControlledAPI.Get("/logout", login.Logout)
			AccessControlledAPI.Post("/user", user.Create)
		}

		//admin.Setup(v1)

		//userApi.Get("/permission")
	}
}
