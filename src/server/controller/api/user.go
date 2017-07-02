package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp-api/src/server/common"
)

type UserController struct{}

func (self *UserController) AuthorizationUser(ctx *iris.Context) {

}

func (self *UserController) GetSessionUser(ctx *iris.Context) {

}

func (self *UserController) GetById(ctx *iris.Context) {
	common.Render(ctx, "12000100", nil)
}
