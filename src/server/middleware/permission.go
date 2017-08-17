package middleware

import (
	"gopkg.in/kataras/iris.v5"
	"strings"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/service/permission"
)

//检验控制器访问权限的中间件
func AccessControlMiddleware(ctx *iris.Context) {
	var (
		userRoleRelService = permission.UserRoleRelService{}
		rolePermissionRelService = permission.RolePermissionRelService{}
		permissionActionRelService = permission.PermissionActionRelService{}
		actionService = permission.ActionService{}
	)
	currentUserID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000001", nil)
		return
	}

	roleIDs, err := userRoleRelService.GetRoleIDsByUserID(currentUserID)
	if err != nil {
		common.Render(ctx, "27010112", nil)
		return
	}
	permissionIDs, err := rolePermissionRelService.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		common.Render(ctx, "27010117", nil)
		return
	}
	if len(permissionIDs) != 0 {
		actionIDs, err := permissionActionRelService.GetActionIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27010115", nil)
			return
		}
		if len(actionIDs) != 0 {
			actionList, err := actionService.GetListByIDs(actionIDs)
			if err != nil {
				common.Render(ctx, "27010116", nil)
				return
			}
			for _, action := range *actionList {
				if strings.EqualFold(action.HandlerName, ctx.GetHandlerName()) {
					ctx.Next()
					return
				}
			}
		}
	}
	common.Render(ctx, "000005", nil)
}
