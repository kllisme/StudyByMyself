package api

import (
	"github.com/bitly/go-simplejson"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp-api/src/server/common"
)

type LoginController struct {
}

func (self *LoginController) Login(ctx *iris.Context) {
	params := simplejson.New()
	err := ctx.ReadJSON(&params)
	if err != nil {
		common.Render(ctx, "12000501", err)
		return
	}
	app := params.Get("app").MustString()
	if app == "WECHAT" {
		self.WechatLogin(ctx)
	} else if app == "MOBILE" {
		self.PassportLogin(ctx)
	} else {
		common.Render(ctx, "12000502", nil)
		return
	}
}

//微信登录
func (self *LoginController) WechatLogin(ctx *iris.Context) {

}

//手机短信验证码，手机密码登录
func (self *LoginController) PassportLogin(ctx *iris.Context) {

}

func (self *LoginController) Logout(ctx *iris.Context) {

}
