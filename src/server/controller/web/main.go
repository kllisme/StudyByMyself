package web

import (
	"gopkg.in/kataras/iris.v5"
)

type WebController struct {
}

func (self *WebController) Home(ctx *iris.Context) {
	ctx.Render("home.html", map[string]interface{}{"title": "Home"})
}

func (self *WebController) Captcha(ctx *iris.Context) {

}
