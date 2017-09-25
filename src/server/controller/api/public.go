package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"time"
)

type PublicController struct{}
//
//func (self *PublicController) Token(ctx *iris.Context) {
//	tokenService := &service.TokenService{}
//	token, err := tokenService.Token(ctx)
//	if err != nil {
//		common.Render(ctx, "12030001", err)
//		return
//	}
//	common.Render(ctx, "12030000", token)
//}

func (self *PublicController) Timestamp(ctx *iris.Context) {
	common.Render(ctx, "04000100", time.Now().Unix())
}
