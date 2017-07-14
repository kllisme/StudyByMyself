package api

import (
	"github.com/bitly/go-simplejson"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	"github.com/spf13/viper"
)

type LoginController struct {
}

func (self *LoginController) Login(ctx *iris.Context) {

	captchaKey := viper.GetString("server.captcha.key")

	userService := service.UserService{}
	tokenService := service.TokenService{}
	//每次调用返回时都清一次图片验证码
	var returnCleanCaptcha = func() {
		ctx.Session().Delete(captchaKey)
	}

	params := simplejson.New()
	err := ctx.ReadJSON(&params)
	if err != nil {
		common.Render(ctx, "27010102", err)
		return
	}
	account := params.Get("account").MustString()
	password := params.Get("password").MustString()
	captcha := params.Get("captcha").MustString()

	/*判断不能为空*/
	if account == "" {
		common.Render(ctx, "27010103", nil)
		returnCleanCaptcha()
		return
	}
	if password == "" {
		common.Render(ctx, "27010104", nil)
		returnCleanCaptcha()
		return
	}
	if captcha == "" {
		common.Render(ctx, "27010105", nil)
		returnCleanCaptcha()
		return
	}

	captchaCache := ctx.Session().GetString(captchaKey)
	if captchaCache == "" {
		common.Render(ctx, "27010106", nil)
		returnCleanCaptcha()
		return
	}

	if captchaCache != captcha {
		common.Render(ctx, "27010107", nil)
		returnCleanCaptcha()
		return
	}

	userEntity, err := userService.GetByAccount(account)
	if err != nil {
		common.Render(ctx, "27010108", err)
		returnCleanCaptcha()
		return
	}
	if userEntity.Password != password {
		common.Render(ctx, "27010109", err)
		returnCleanCaptcha()
		return
	}
	ctx.Session().Set(viper.GetString("server.session.user.id"), userEntity.ID)
	token, err := tokenService.Token(ctx)
	if err != nil {
		common.Render(ctx, "27010110", err)
		returnCleanCaptcha()
		return
	}
	common.Render(ctx, "12030000", token)
	returnCleanCaptcha()
}

//微信登录
func (self *LoginController) WechatLogin(ctx *iris.Context) {

}

//手机短信验证码，手机密码登录
func (self *LoginController) PassportLogin(ctx *iris.Context) {

}

func (self *LoginController) Logout(ctx *iris.Context) {

}
