package middleware

import (
	"gopkg.in/kataras/iris.v5"
)

//检验控制器访问权限的中间件
func AccessControlMiddleware(ctx *iris.Context) {
	//logrus.Debug(ctx.GetHandlerName())
	//info := payload.SessionInfo{}
	//jsonString := ctx.Session().GetString(viper.GetString("server.session.user.key"))
	//if err := json.Unmarshal([]byte(jsonString), &info); err != nil {
	//	common.Render(ctx, "000001", nil)
	//}
	//for _, action := range *info.ActionList {
	//	logrus.Debug(action.HandlerName)
	//	if strings.EqualFold(action.HandlerName, ctx.GetHandlerName()) {
	ctx.Next()
	//		return
	//	}
	//}
	//common.Render(ctx, "000005", nil)
}
