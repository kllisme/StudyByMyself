package api

import (
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	"strings"
	"github.com/bitly/go-simplejson"
)

type LoginController struct {
}

func (self *LoginController) Login(ctx *iris.Context) {
	var (
		captchaKey = viper.GetString("server.captcha.key")
		userService = service.UserService{}
		tokenService = service.TokenService{}
	)

	//每次调用返回时都清一次图片验证码
	defer ctx.Session().Delete(captchaKey)

	params := simplejson.New()
	err := ctx.ReadJSON(&params)
	if err != nil {
		common.Render(ctx, "27010102", nil)
		return
	}
	account := strings.TrimSpace(params.Get("account").MustString())
	password := strings.TrimSpace(params.Get("password").MustString())
	captcha := strings.TrimSpace(params.Get("captcha").MustString())

	/*判断不能为空*/
	if account == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	if password == "" {
		common.Render(ctx, "27010104", nil)
		return
	}
	if captcha == "" {
		common.Render(ctx, "27010105", nil)
		return
	}

	captchaCache := ctx.Session().GetString(captchaKey)
	if captchaCache == "" {
		common.Render(ctx, "27010106", nil)
		return
	}

	if captchaCache != captcha {
		common.Render(ctx, "27010107", nil)
		return
	}

	userEntity, err := userService.GetByAccount(account)
	if err != nil {
		common.Render(ctx, "27010108", nil)
		return
	}
	if userEntity.Password != password {
		common.Render(ctx, "27010109", nil)
		return
	}

	token, err := tokenService.Token(ctx)
	if err != nil {
		common.Render(ctx, "27010110", err)
		return
	}
	ctx.Session().Set(viper.GetString("server.session.user.id"), userEntity.ID)
	common.Render(ctx, "27010100", token)
}

//微信登录
func (self *LoginController) WechatLogin(ctx *iris.Context) {

}

//手机短信验证码，手机密码登录
func (self *LoginController) PassportLogin(ctx *iris.Context) {

}

func (self *LoginController) Logout(ctx *iris.Context) {
	ctx.SessionDestroy()
	common.Render(ctx, "27010200", nil)
	return
}
