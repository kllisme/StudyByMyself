package admin

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api"
	"maizuo.com/soda/erp/api/src/server/controller/api/permission"
)

func Setup(v iris.MuxAPI) {
	var (
		userCtrl = &api.UserController{}
		roleCtrl = &permission.RoleController{}
		menuCtrl = &permission.MenuController{}
		elementCtrl = &permission.ElementController{}
		actionCtrl = &permission.ActionController{}
		permissionCtrl = &permission.PermissionController{}
	)
	_api := v.Party("/")
	_api.Post("/users", userCtrl.Create)
	_api.Get("/users/:id", userCtrl.GetByID)
	_api.Get("/users", userCtrl.Paging)
	_api.Put("/users/:id", userCtrl.Update)
	_api.Delete("/users/:id", userCtrl.Delete)
	_api.Put("/users/:id/password", userCtrl.ResetPassword)
	_api.Put("/users/:id/roles", userCtrl.AssignRoles)
	_api.Get("/users/:id/roles", userCtrl.GetRoles)


	_api.Post("/roles", roleCtrl.Create)
	_api.Get("/roles", roleCtrl.GetAll)
	_api.Delete("/roles/:id", roleCtrl.Delete)
	_api.Get("/roles/:id", roleCtrl.GetByID)
	_api.Put("/roles/:id", roleCtrl.Update)
	_api.Put("/roles/:id/permissions", roleCtrl.AssignPermissions)
	_api.Get("/roles/:id/permissions", roleCtrl.GetPermissions)

	_api.Post("/menus", menuCtrl.Create)
	_api.Delete("/menus/:id", menuCtrl.Delete)
	_api.Put("/menus/:id", menuCtrl.Update)
	_api.Get("/menus/:id", menuCtrl.GetByID)
	_api.Get("/menus", menuCtrl.Paging)

	_api.Post("/permissions", permissionCtrl.Create)
	_api.Delete("/permissions/:id", permissionCtrl.Delete)
	_api.Put("/permissions/:id", permissionCtrl.Update)
	_api.Get("/permissions/:id", permissionCtrl.GetByID)
	_api.Get("/permissions", permissionCtrl.Paging)
	_api.Get("/permissions/:id/menus", permissionCtrl.GetMenus)
	_api.Put("/permissions/:id/menus", permissionCtrl.AssignMenus)
	_api.Get("/permissions/:id/elements", permissionCtrl.GetElements)
	_api.Put("/permissions/:id/elements", permissionCtrl.AssignElements)
	_api.Get("/permissions/:id/actions", permissionCtrl.GetActions)
	_api.Put("/permissions/:id/actions", permissionCtrl.AssignActions)

	_api.Post("/actions", actionCtrl.Create)
	_api.Delete("/actions/:id", actionCtrl.Delete)
	_api.Put("/actions/:id", actionCtrl.Update)
	_api.Get("/actions/:id", actionCtrl.GetByID)
	_api.Get("/actions", actionCtrl.Paging)

	_api.Post("/elements", elementCtrl.Create)
	_api.Delete("/elements/:id", elementCtrl.Delete)
	_api.Put("/elements/:id", elementCtrl.Update)
	_api.Get("/elements/:id", elementCtrl.GetByID)
	_api.Get("/elements", elementCtrl.Paging)

}
