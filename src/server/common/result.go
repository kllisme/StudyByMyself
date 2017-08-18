package common

import (
	"fmt"

	"maizuo.com/soda/erp/api/src/server/entity"

	"github.com/go-errors/errors"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
)

type Result struct {
	Status      string      `json:"status"`  //UNLOGIN
	Data        interface{} `json:"data"`    //{}
	Msg         string      `json:"message"` //中文:系统开小差了，请稍后再试！
	Description string      `json:"-"`       //异常具体描述,如：解析JSON异常
	Exception   interface{} `json:"-"`       //{}
	Code        string      `json:"code"`    //系统唯一错误码：101111111
	IsError     bool        `json:"-"`
}

func (self *Result) New(code string, data interface{}) *Result {
	prefix := "status.biz."
	msg := viper.GetString(prefix + code + "." + "msg")
	description := viper.GetString(prefix + code + "." + "description")
	isError := viper.GetBool(prefix + code + "." + "isError")
	status := viper.GetString(prefix + code + "." + "status")
	result := &Result{
		Status:      status,
		Msg:         msg,
		Description: description,
		Code:        code,
		IsError:     isError,
	}
	switch data.(type) {
	case *Result:
		result.Data = struct{}{}
	case *entity.Pagination:
		result.Data = struct{}{}
	default:
		result.Data = data
	}
	if result.IsError {
		errorStack := fmt.Sprintf("%v\n", errors.Wrap(data, 1).ErrorStack())
		if len(errorStack) > 500 {
			errorStack = errorStack[:500]
		}
		result.Exception = errorStack
		result.Data = struct{}{}
	} else {
		result.Exception = ""
	}
	return result
}

func Error(code string, data interface{}) *Result {
	if code == "" {
		code = "000001"
	}
	prefix := "status.service."
	msg := viper.GetString(prefix + code + "." + "msg")
	description := viper.GetString(prefix + code + "." + "description")
	result := &Result{
		Status:      code,
		Msg:         msg,
		Description: description,
		Code:        code,
	}
	if data == nil {
		errorStack := fmt.Sprintf("%v\n", errors.Wrap(data, 1).ErrorStack())
		if len(errorStack) > 500 {
			errorStack = errorStack[:500]
		}
		result.Exception = errorStack
	}
	return result
}

func Render(ctx *iris.Context, code string, data interface{}) {
	result := &Result{}
	_result := result.New(code, data)
	ctx.JSON(iris.StatusOK, _result)
	Log(ctx, _result)
}
