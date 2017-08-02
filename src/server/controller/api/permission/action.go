package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type ActionController struct {

}

func (self *ActionController)GetByID(ctx *iris.Context) {
	actionService := permission.ActionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
	}
	action, err := actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040100", action)
}

func (self *ActionController)Query(ctx *iris.Context) {
	actionService := permission.ActionService{}
	conditions := model.Action{}

	method := strings.TrimSpace(ctx.URLParam("method"))
	handlerName := ctx.URLParam("handler_name")
	//TODO complete all query conditions
	if method != "" {
		conditions.Method = method
	}
	if handlerName != "" {
		conditions.HandlerName = handlerName
	}
	actionList, err := actionService.Query(&conditions)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27040200", actionList)
}

func (self *ActionController)Create(ctx *iris.Context) {
	actionService := permission.ActionService{}
	action := model.Action{}
	if err := ctx.ReadJSON(&action); err != nil {
		common.Render(ctx, "27040301", nil)
	}
	if err := actionService.Create(&action); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040300", action)
}

func (self *ActionController)Delete(ctx *iris.Context) {
	actionService := permission.ActionService{}
	id, err := ctx.URLParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
	}
	if err := actionService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040400", nil)
}

func (self *ActionController)Update(ctx *iris.Context) {
	actionService := permission.ActionService{}
	params := simplejson.New()

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
	}
	action, err := actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
	}
	err = ctx.ReadJSON(&params)
	if err != nil {
		common.Render(ctx, "27040501", nil)
		return
	}
	action.API = strings.TrimSpace(params.Get("api").MustString())
	action.Description = strings.TrimSpace(params.Get("description").MustString())
	action.HandlerName = strings.TrimSpace(params.Get("handlerName").MustString())
	action.Method = strings.TrimSpace(params.Get("method").MustString())

	result, err := actionService.Update(action);
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040500", result)
}
