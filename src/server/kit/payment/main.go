package payment

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
)

func Render(ctx *iris.Context, code string, data interface{}) {
	result := &common.Result{}
	_result := result.New(code, data)
	common.Log(ctx, _result)
}
