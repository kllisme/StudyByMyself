package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/middleware"
	"maizuo.com/soda/erp/api/src/server/controller/api/permission"
	"maizuo.com/soda/erp/api/src/server/route/api/admin"
)

func Api(app *iris.Framework) {

	var (
		userCtrl = &api.UserController{}
		//public = &api.PublicController{}
		captchaCtrl = &api.CaptchaController{}
		loginCtrl = &api.LoginController{}
		roleCtrl = &permission.RoleController{}
		//menuCtrl = &permission.MenuController{}
		//actionCtrl = &permission.ActionController{}
		//permissionCtrl = &permission.PermissionController{}
	)
	v1 := app.Party("/v1", func(ctx *iris.Context) {
		ctx.Next()
	})
	{
		v1.Get("/captcha.png", captchaCtrl.Captcha)

		//为跨域请求设定入口
		v1.UseFunc(common.CORS.Serve)
		v1.Options("/*anything", common.CORS.Serve)

		//v1.Get("/token", public.Token)
		v1.Post("/login", loginCtrl.Login)

		v1.UseFunc(common.Authorization)
		v1.Get("/session/info", userCtrl.GetSessionInfo)
		v1.Get("/session/user", userCtrl.GetSessionInfo)

		//控制访问权限的接口
		accessControlledAPI := v1.UseFunc(middleware.AccessControlMiddleware)
		{
			accessControlledAPI.Get("/logout", loginCtrl.Logout)
			accessControlledAPI.Post("/user", userCtrl.Create)
			accessControlledAPI.Get("/user", userCtrl.GetById)
			accessControlledAPI.Put("/user", userCtrl.Update)
			accessControlledAPI.Put("/user/password", userCtrl.ResetPassword)
			//accessControlledAPI.Post("/user/role", user.AssignRoles)
			accessControlledAPI.Put("/user/role", userCtrl.ChangeRoles)

			admin.Setup(accessControlledAPI)
			permissionAPI := accessControlledAPI.Party("/")
			{
				//permissionAPI.Post("/role", roleCtrl.Create)
				//permissionAPI.Delete("/role", roleCtrl.Delete)
				//permissionAPI.Put("/role", roleCtrl.Update)
				//permissionAPI.Get("/role", roleCtrl.GetByID)
				//permissionAPI.Get("/roles", roleCtrl.GetAll)
				//
				//permissionAPI.Post("/menu", menuCtrl.Create)
				//permissionAPI.Delete("/menu", menuCtrl.Delete)
				//permissionAPI.Put("/menu", menuCtrl.Update)
				//permissionAPI.Get("/menu", menuCtrl.GetByID)
				//permissionAPI.Get("/menus", menuCtrl.GetAll)
				//
				//permissionAPI.Post("/action", actionCtrl.Create)
				//permissionAPI.Delete("/action", actionCtrl.Delete)
				//permissionAPI.Put("/action", actionCtrl.Update)
				//permissionAPI.Get("/action", actionCtrl.GetByID)
				//permissionAPI.Get("/actions", actionCtrl.GetAll)
				//
				//permissionAPI.Post("/permission", permissionCtrl.Create)
				//permissionAPI.Delete("/permission", permissionCtrl.Delete)
				//permissionAPI.Put("/permission", permissionCtrl.Update)
				//permissionAPI.Get("/permission", permissionCtrl.GetByID)
				//permissionAPI.Get("/permissions", permissionCtrl.GetAll)

				permissionAPI.Post("/menu", roleCtrl.AssignPermissions)
			}
		}

		//admin.Setup(v1)

		//userApi.Get("/permission")
	}
}
