package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/middleware"
	"maizuo.com/soda/erp/api/src/server/route/api/admin"
)

func Api(app *iris.Framework) {

	var (
		userCtrl = &api.UserController{}
		//public = &api.PublicController{}
		captchaCtrl = &api.CaptchaController{}
		loginCtrl   = &api.LoginController{}

		billCtrl      = &api.BillController{}
		dailyBillCtrl = &api.DailyBillController{}
	)


	v1 := app.Party("/v1", func(ctx *iris.Context) {
		ctx.Next()
	})

	{

		v1.Post("/settlement/actions/wechatPay", billCtrl.WechatPay)
		v1.Post("/settlement/alipay/notification",billCtrl.AlipayNotification)

		v1.Get("/captcha.png", captchaCtrl.Captcha)

		//为跨域请求设定入口
		v1.UseFunc(common.CORS.Serve)
		v1.Options("/*anything", common.CORS.Serve)

		//v1.Get("/token", public.Token)
		v1.Post("/login", loginCtrl.Login)


		//jwt校验
		v1.UseFunc(common.Authorization)
		v1.Get("/profile/session", userCtrl.GetSessionInfo)

		//v1.Get("/profile/user", userCtrl.GetSessionInfo)

		api := v1.Party("", func(ctx *iris.Context) {
			ctx.Next()
		})
		{
			//api.UseFunc(middleware.BillRoleControlMiddleware)
			api.Get("/bills", billCtrl.ListByAccountType)
			api.Get("/bills/:id", dailyBillCtrl.ListByBillId)

			api.Get("/daily-bills/:id", dailyBillCtrl.DetailsById)

			api.Post("/settlement/actions/pay", billCtrl.BatchPay)


			//api.Post("/settlement/actions/cancel", billCtrl.CancelBatchAliPay)
		}

		//控制访问权限的接口
		accessControlledAPI := v1.UseFunc(middleware.AccessControlMiddleware)
		{
			accessControlledAPI.Get("/logout", loginCtrl.Logout)
			accessControlledAPI.Put("/profile/password", userCtrl.ChangePassword)
			//本系统账号及权限管理
			admin.Setup(accessControlledAPI)
		}
	}
}
