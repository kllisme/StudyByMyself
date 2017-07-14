package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	"github.com/spf13/viper"
)

type UserController struct{}

func (self *UserController) AuthorizationUser(ctx *iris.Context) {

}

func (self *UserController) GetSessionUser(ctx *iris.Context) {
	userService := service.UserService{}
	currentUserId,_ := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	userEntity,_:= userService.GetById(currentUserId)
	common.Render(ctx,"12010000",userEntity)
}

func (self *UserController) GetById(ctx *iris.Context) {
	userService := service.UserService{}
	id,_ := ctx.URLParamInt("id")
	userEntity,_:= userService.GetById(id)
	common.Render(ctx,"12010000",userEntity)
}
