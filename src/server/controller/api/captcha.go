package api

import (
	"image/color"
	"image/png"

	"github.com/afocus/captcha"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
)

type CaptchaController struct {
}

func (self *CaptchaController) Captcha(ctx *iris.Context) {
	_cap := captcha.New()
	fontPath := viper.GetString("server.captcha.font.path")
	captchaKey := viper.GetString("server.captcha.key")
	if err := _cap.SetFont(fontPath); err != nil {
		panic(err.Error())
	}
	_cap.SetSize(128, 64)
	_cap.SetDisturbance(captcha.MEDIUM)
	_cap.SetFrontColor(color.RGBA{0, 0, 0, 255})
	_cap.SetBkgColor(color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0})
	img, _captcha := _cap.Create(4, captcha.NUM)
	png.Encode(ctx.Response.BodyWriter(), img)
	ctx.Session().Set(captchaKey, _captcha)
}
