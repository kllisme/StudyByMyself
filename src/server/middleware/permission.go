package middleware

import (
	"gopkg.in/kataras/iris.v5"
	"github.com/Sirupsen/logrus"
	"maizuo.com/soda/erp/api/src/server/payload"
	"strings"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/spf13/viper"
	"encoding/json"
)

//检验控制器访问权限的中间件
func AccessControlMiddleware(ctx *iris.Context) {
	info := payload.SessionInfo{}
	jsonString := ctx.Session().GetString(viper.GetString("server.session.user.key"))
	if err := json.Unmarshal([]byte(jsonString), &info); err != nil {
		common.Render(ctx, "000001", nil)
	}
	for _, action := range *info.ActionList {
		logrus.Debug(action.HandlerName)
		if strings.EqualFold(action.HandlerName, ctx.GetHandlerName()) {
			ctx.Next()
			return
		}
	}
	common.Render(ctx, "000005", nil)
}
