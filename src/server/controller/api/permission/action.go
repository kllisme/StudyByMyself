package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
)

type ActionController struct {

}

func (self *ActionController)GetByID(ctx *iris.Context) {
	actionService := permission.ActionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
	}
	action, err := actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27040100", action)
}

func (self *ActionController)Paging(ctx *iris.Context) {
	actionService := permission.ActionService{}
	method := strings.TrimSpace(ctx.URLParam("method"))
	handlerName := strings.TrimSpace(ctx.URLParam("handler_name"))
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := actionService.Paging(page,perPage,handlerName,method)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27040200", result)
}

func (self *ActionController)Create(ctx *iris.Context) {
	actionService := permission.ActionService{}
	action := model.Action{}
	if err := ctx.ReadJSON(&action); err != nil {
		common.Render(ctx, "27040301", err)
	}
	if err := actionService.Create(&action); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27040300", action)
}

func (self *ActionController)Delete(ctx *iris.Context) {
	actionService := permission.ActionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
	}
	if err := actionService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27040400", nil)
}

func (self *ActionController)Update(ctx *iris.Context) {
	actionService := permission.ActionService{}
	action := model.Action{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
	}
	_, err = actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
	}
	err = ctx.ReadJSON(&action)
	if err != nil {
		common.Render(ctx, "27040501", err)
		return
	}

	action.API = strings.TrimSpace(action.API)
	action.Description = strings.TrimSpace(action.Description)
	action.HandlerName = strings.TrimSpace(action.HandlerName)
	action.Method = strings.TrimSpace(action.Method)
	action.ID = id
	result, err := actionService.Update(&action)
	if err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27040500", result)
}
