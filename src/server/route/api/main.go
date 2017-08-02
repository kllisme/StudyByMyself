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
		menuCtrl = &permission.MenuController{}
		elementCtrl = &permission.ElementController{}
		actionCtrl = &permission.ActionController{}
		permissionCtrl = &permission.PermissionController{}
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
		v1.Get("/profile/session", userCtrl.GetSessionInfo)

		//v1.Get("/profile/user", userCtrl.GetSessionInfo)

		//控制访问权限的接口
		accessControlledAPI := v1.UseFunc(middleware.AccessControlMiddleware)
		{
			accessControlledAPI.Get("/logout", loginCtrl.Logout)

			accessControlledAPI.Post("/user", userCtrl.Create)
			accessControlledAPI.Get("/user/:id", userCtrl.GetById)
			accessControlledAPI.Get("/users", userCtrl.Paging)

			accessControlledAPI.Put("/user/:id", userCtrl.Update)
			accessControlledAPI.Delete("/user/:id", userCtrl.Delete)
			accessControlledAPI.Put("/user/:id/password", userCtrl.ResetPassword)
			accessControlledAPI.Put("/user/:id/roles", userCtrl.AssignRoles)
			accessControlledAPI.Get("/user/:id/roles", userCtrl.GetRoles)

			accessControlledAPI.Put("/profile/password", userCtrl.ChangePassword)


			accessControlledAPI.Post("/role", roleCtrl.Create)
			accessControlledAPI.Get("/roles", roleCtrl.GetAll)
			accessControlledAPI.Delete("/role/:id", roleCtrl.Delete)
			accessControlledAPI.Get("/role/:id", roleCtrl.GetByID)
			accessControlledAPI.Put("/role/:id",roleCtrl.Update)
			accessControlledAPI.Put("/role/:id/permissions", roleCtrl.AssignPermissions)
			accessControlledAPI.Get("/role/:id/permissions", roleCtrl.GetPermissions)

			admin.Setup(accessControlledAPI)
			permissionAPI := accessControlledAPI.Party("/")
			{

				permissionAPI.Post("/menu", menuCtrl.Create)
				permissionAPI.Delete("/menu/:id", menuCtrl.Delete)
				permissionAPI.Put("/menu/:id", menuCtrl.Update)
				permissionAPI.Get("/menu/:id", menuCtrl.GetByID)
				permissionAPI.Get("/menus", menuCtrl.Paging)

				permissionAPI.Post("/permission", permissionCtrl.Create)
				permissionAPI.Delete("/permission/:id", permissionCtrl.Delete)
				permissionAPI.Put("/permission/:id", permissionCtrl.Update)
				permissionAPI.Get("/permission/:id", permissionCtrl.GetByID)
				permissionAPI.Get("/permissions", permissionCtrl.Paging)
				permissionAPI.Get("/permission/:id/menu", permissionCtrl.GetMenus)
				permissionAPI.Put("/permission/:id/menu", permissionCtrl.AssignMenus)
				permissionAPI.Get("/permission/:id/element", permissionCtrl.GetElements)
				permissionAPI.Put("/permission/:id/element", permissionCtrl.AssignElements)
				permissionAPI.Get("/permission/:id/action", permissionCtrl.GetActions)
				permissionAPI.Put("/permission/:id/action", permissionCtrl.AssignActions)

				permissionAPI.Post("/action", actionCtrl.Create)
				permissionAPI.Delete("/action/:id", actionCtrl.Delete)
				permissionAPI.Put("/action/:id", actionCtrl.Update)
				permissionAPI.Get("/action/:id", actionCtrl.GetByID)
				permissionAPI.Get("/actions", actionCtrl.Query)

				permissionAPI.Post("/element", elementCtrl.Create)
				permissionAPI.Delete("/element/:id", elementCtrl.Delete)
				permissionAPI.Put("/element/:id", elementCtrl.Update)
				permissionAPI.Get("/element/:id", elementCtrl.GetByID)
				permissionAPI.Get("/elements", elementCtrl.Paging)

			}
		}

		//admin.Setup(v1)

		//userApi.Get("/permission")
	}
}
