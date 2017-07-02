package admin

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp-api/src/server/controller/api"
)

func Setup(v iris.MuxAPI) {
	var (
		user = &api.UserController{}
	)
	_api := v.Party("/admin")
	_api.Get("/a", user.GetById)
	_api.Get("/detail", user.GetById)

}
