package middleware

import (
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"maizuo.com/soda/erp/api/src/server/service/permission"
)

//检验控制器访问权限的中间件
func AccessControlMiddleware(ctx *iris.Context) {
	var (
		userRoleRelService         = permission.UserRoleRelService{}
		rolePermissionRelService   = permission.RolePermissionRelService{}
		permissionActionRelService = permission.PermissionActionRelService{}
		actionService              = permission.ActionService{}
	)
	currentUserID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000001", nil)
		return
	}
	roleIDs, err := userRoleRelService.GetRoleIDsByUserID(currentUserID)
	if err != nil {
		common.Render(ctx, "000005", nil)
		return
	}
	permissionIDs, err := rolePermissionRelService.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		common.Render(ctx, "000005", nil)
		return
	}
	if len(permissionIDs) == 0 {
		common.Render(ctx, "000005", nil)
		return
	}
	actionIDs, err := permissionActionRelService.GetActionIDsByPermissionIDs(permissionIDs)
	if err != nil {
		common.Render(ctx, "000005", nil)
		return
	}
	if len(actionIDs) == 0 {
		common.Render(ctx, "000005", nil)
		return
	}
	actionList, err := actionService.GetListByIDs(actionIDs)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	hasAuthorized := false
	handlerName := functions.ExtractHandlerName(ctx.GetHandlerName())
	for _, action := range *actionList {
		if strings.EqualFold(action.HandlerName, handlerName) {
			hasAuthorized = true
			break
		}
	}
	if !hasAuthorized {
		common.Render(ctx, "000005", nil)
		return
	}
	ctx.Next()
}
