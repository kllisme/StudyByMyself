package middleware

import (
	"gopkg.in/kataras/iris.v5"
	"github.com/spf13/viper"
	"encoding/json"
	"maizuo.com/soda/erp/api/src/server/payload"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/Sirupsen/logrus"
)

//检验控制器访问权限的中间件
func AccessControlMiddleware(ctx *iris.Context) {
	//runtime.FuncForPC(reflect.ValueOf(ctx.Middleware[len(ctx.Middleware)-1]).Pointer()).Name()
	logrus.Debug(ctx.GetHandlerName())
	info := payload.SessionInfo{}
	jsonString := ctx.Session().GetString(viper.GetString("server.session.user.key"))
	if err := json.Unmarshal([]byte(jsonString), &info); err != nil {
		common.Render(ctx, "000001", nil)
	}
	for _, api := range *info.APIList {
		logrus.Debug(api.Name)
		if strings.EqualFold(api.Name, ctx.GetHandlerName()) {
			ctx.Next()
			return
		}
	}
	common.Render(ctx, "000005", nil)
}
