package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/middleware"
	"maizuo.com/soda/erp/api/src/server/controller/api/permission"
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

			accessControlledAPI.Post("/users", userCtrl.Create)
			accessControlledAPI.Get("/users/:id", userCtrl.GetById)
			accessControlledAPI.Get("/users", userCtrl.Paging)

			accessControlledAPI.Put("/users/:id", userCtrl.Update)
			accessControlledAPI.Delete("/users/:id", userCtrl.Delete)
			accessControlledAPI.Put("/users/:id/password", userCtrl.ResetPassword)
			accessControlledAPI.Put("/users/:id/roles", userCtrl.AssignRoles)
			accessControlledAPI.Get("/users/:id/roles", userCtrl.GetRoles)

			accessControlledAPI.Put("/profile/password", userCtrl.ChangePassword)

			accessControlledAPI.Post("/roles", roleCtrl.Create)
			accessControlledAPI.Get("/roles", roleCtrl.GetAll)
			accessControlledAPI.Delete("/roles/:id", roleCtrl.Delete)
			accessControlledAPI.Get("/roles/:id", roleCtrl.GetByID)
			accessControlledAPI.Put("/roles/:id", roleCtrl.Update)
			accessControlledAPI.Put("/roles/:id/permissions", roleCtrl.AssignPermissions)
			accessControlledAPI.Get("/roles/:id/permissions", roleCtrl.GetPermissions)

			//admin.Setup(accessControlledAPI)
			permissionAPI := accessControlledAPI.Party("/")
			{

				permissionAPI.Post("/menus", menuCtrl.Create)
				permissionAPI.Delete("/menus/:id", menuCtrl.Delete)
				permissionAPI.Put("/menus/:id", menuCtrl.Update)
				permissionAPI.Get("/menus/:id", menuCtrl.GetByID)
				permissionAPI.Get("/menus", menuCtrl.Paging)

				permissionAPI.Post("/permissions", permissionCtrl.Create)
				permissionAPI.Delete("/permissions/:id", permissionCtrl.Delete)
				permissionAPI.Put("/permissions/:id", permissionCtrl.Update)

				permissionAPI.Get("/permissions/:id", permissionCtrl.GetByID)
				permissionAPI.Get("/permissions", permissionCtrl.Paging)

				permissionAPI.Get("/permissions/:id/menus", permissionCtrl.GetMenus)
				permissionAPI.Put("/permissions/:id/menus", permissionCtrl.AssignMenus)
				permissionAPI.Get("/permissions/:id/elements", permissionCtrl.GetElements)
				permissionAPI.Put("/permissions/:id/elements", permissionCtrl.AssignElements)
				permissionAPI.Get("/permissions/:id/actions", permissionCtrl.GetActions)
				permissionAPI.Put("/permissions/:id/actions", permissionCtrl.AssignActions)

				permissionAPI.Post("/actions", actionCtrl.Create)
				permissionAPI.Delete("/actions/:id", actionCtrl.Delete)
				permissionAPI.Put("/actions/:id", actionCtrl.Update)
				permissionAPI.Get("/actions/:id", actionCtrl.GetByID)
				permissionAPI.Get("/actions", actionCtrl.Query)

				permissionAPI.Post("/elements", elementCtrl.Create)
				permissionAPI.Delete("/elements/:id", elementCtrl.Delete)
				permissionAPI.Put("/elements/:id", elementCtrl.Update)
				permissionAPI.Get("/elements/:id", elementCtrl.GetByID)
				permissionAPI.Get("/elements", elementCtrl.Paging)

			}
		}

		//admin.Setup(v1)

		//userApi.Get("/permission")
	}
}
