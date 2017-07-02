package common

import (
	"fmt"

	"github.com/go-errors/errors"
	"gopkg.in/kataras/iris.v5"
)

type Result struct {
	Status      string      `json:"status"`  //UNLOGIN
	Data        interface{} `json:"data"`    //{}
	Msg         string      `json:"message"` //中文:系统开小差了，请稍后再试！
	Description string      `json:"-"`       //异常具体描述,如：解析JSON异常
	Exception   interface{} `json:"-"`       //{}
	Code        string      `json:"-"`       //系统唯一错误码：101111111
	IsError     bool        `json:"-"`
}

func (self *Result) New(code string, data interface{}) *Result {
	prefix := "biz."
	msg := StatusConfig.GetString(prefix + code + "." + "msg")
	description := StatusConfig.GetString(prefix + code + "." + "description")
	isError := StatusConfig.GetBool(prefix + code + "." + "isError")
	status := StatusConfig.GetString(prefix + code + "." + "status")
	result := &Result{
		Status:      status,
		Msg:         msg,
		Description: description,
		Code:        code,
		IsError:     isError,
	}
	switch data.(type) {
	case *Result:
		result.Data = nil
	default:
		result.Data = data
	}
	if result.IsError {
		result.Exception = fmt.Sprintf("%v\n", errors.Wrap(data, 1).ErrorStack())[:500]
		result.Data = nil
	} else {
		result.Exception = ""
	}
	return result
}

func Error(code string, data interface{}) *Result {
	if code == "" {
		code = "000001"
	}
	prefix := "service."
	msg := StatusConfig.GetString(prefix + code + "." + "msg")
	description := StatusConfig.GetString(prefix + code + "." + "description")
	result := &Result{
		Status:      code,
		Msg:         msg,
		Description: description,
		Code:        code,
	}
	if data == nil {
		result.Exception = fmt.Sprintf("%v\n", errors.Wrap(data, 1).ErrorStack())[:500]
	}
	return result
}

func Render(ctx *iris.Context, code string, data interface{}) {
	result := &Result{}
	_result := result.New(code, data)
	Log(ctx, _result)
	ctx.JSON(iris.StatusOK, _result)
}
