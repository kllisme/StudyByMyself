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
		common.Render(ctx, "000003", nil)
	}
	action, err := actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "", action)
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
	if handlerName != ""{
		conditions.HandlerName = handlerName
	}
	actionList, err := actionService.Query(&conditions)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27030200", actionList)
}

func (self *ActionController)Create(ctx *iris.Context) {
	actionService := permission.ActionService{}
	action := model.Action{}
	if err := ctx.ReadJSON(&action); err != nil {
		common.Render(ctx, "27030301", nil)
	}
	if err := actionService.Create(&action); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27030300", action)
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
	common.Render(ctx, "27030400", nil)
}

func (self *ActionController)Update(ctx *iris.Context) {
	actionService := permission.ActionService{}
	action := model.Action{}

	id, err := ctx.URLParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
	}

	if err := ctx.ReadJSON(&action); err != nil {
		common.Render(ctx, "27030501", nil)
	}

	action.ID = id
	if err := actionService.Update(&action); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27030500", action)
}
