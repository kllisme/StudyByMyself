package middleware

import (
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

//控制财务角色访问权限的中间件
func BillRoleControlMiddleware(ctx *iris.Context) {
	// 或许可以放到中间件
	userRoleService := &permission.UserRoleRelService{}
	userId, _ := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	userRoleList, err := userRoleService.GetRoleIDsByUserID(userId)
	if err != nil {
		common.Logger.Debugln("获取当前操作用户角色失败 userId------", userId)
		common.Render(ctx, "27021101", err)
		return
	}
	// 判断是不是财务或者系统管理员,不是财务的不放行
	if functions.FindIndex(userRoleList, 3) == -1 && functions.FindIndex(userRoleList, 5) == -1 {
		common.Logger.Debugln("获取当前操作用户不具有权限 userRoleList-----", userRoleList)
		common.Render(ctx, "27050801", nil)
		return
	}
}


