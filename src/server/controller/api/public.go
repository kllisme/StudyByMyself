package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
)

type PublicController struct{}

func (self *PublicController) Token(ctx *iris.Context) {
	tokenService := &service.TokenService{}
	token, err := tokenService.Token(ctx)
	if err != nil {
		common.Render(ctx, "12030001", err)
		return
	}
	common.Render(ctx, "12030000", token)
}
