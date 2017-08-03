package server

import (
	"time"

	"github.com/spf13/viper"
	"gopkg.in/iris-contrib/middleware.v5/logger"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/route/api"
)

func SetUpServer() {

	isDevelopment := viper.GetBool("isDevelopment")

	iris.Config.IsDevelopment = isDevelopment

	app := iris.Default //如果使用iris.New(),之前设置的iris.Config无法生效

	app.Use(common.NewRecover())

	if isDevelopment {
		iris.Use(logger.New())
	}

	app.UseFunc(func(ctx *iris.Context) {
		startAt := time.Now().UnixNano() / 1000000
		ctx.Set("startAt", startAt)
		ctx.SetHeader("X-Powered-By", "soda-erp-api/"+viper.GetString("version"))
		ctx.Next()
	})

	app.UseSessionDB(common.SESSION_DB)

	api.Api(app)

	app.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		result := &common.Result{}
		_result := result.New("000003", nil)
		ctx.JSON(iris.StatusOK, _result)
	})

	app.OnError(iris.StatusTooManyRequests, func(ctx *iris.Context) {
		result := &common.Result{}
		_result := result.New("000004", nil)
		ctx.JSON(iris.StatusOK, _result)
	})

	app.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		common.Render(ctx, "000002", nil)
	})

	app.Listen(viper.GetString("server.host") + ":" + viper.GetString("server.port"))

}
