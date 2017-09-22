package api

import (
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/controller/api/finance"
	"maizuo.com/soda/erp/api/src/server/middleware"
	adminApi "maizuo.com/soda/erp/api/src/server/route/api/admin"
	financeApi "maizuo.com/soda/erp/api/src/server/route/api/finance"
	twoApi "maizuo.com/soda/erp/api/src/server/route/api/two"
	publicApi "maizuo.com/soda/erp/api/src/server/route/api/public"
)

func Api(app *iris.Framework) {

	var (
		userCtrl    = &api.UserController{}
		captchaCtrl = &api.CaptchaController{}
		loginCtrl   = &api.LoginController{}
		billCtrl    = &finance.BillController{}
	)

	v1 := app.Party("/v1", func(ctx *iris.Context) {
		ctx.Next()
	})

	v1.Post("/settlement/actions/wechatPay", billCtrl.WechatPay)
	v1.Post("/settlement/alipay/notification", billCtrl.AlipayNotification)

	v1.Get("/captcha.png", captchaCtrl.Captcha)
	v1.StaticFS(viper.GetString("export.loadsPath"), "."+viper.GetString("export.loadsPath"), 2)
	//为跨域请求设定入口
	v1.UseFunc(common.CORS.Serve)
	v1.Options("/*anything", common.CORS.Serve)

	v1.Post("/login", loginCtrl.Login)
	v1.Post("/logout", loginCtrl.Logout)

	//jwt校验
	v1.UseFunc(common.Authorization)

	v1.Get("/profile", userCtrl.GetProfile)

	//控制访问权限的接口
	v1.UseFunc(middleware.AccessControlMiddleware)

	v1.Put("/profile/password", userCtrl.ChangePassword)

	adminApi.Setup(v1)

	financeApi.Setup(v1)

	twoApi.Setup(v1)

	publicApi.Setup(v1)

}
