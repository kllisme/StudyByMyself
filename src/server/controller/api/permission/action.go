package permission

import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
)

type ActionController struct {

}

func (self *ActionController)GetByID(ctx *iris.Context) {
	actionService := mngService.ActionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	action, err := actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "01030100", action)
}

func (self *ActionController)Paging(ctx *iris.Context) {
	actionService := mngService.ActionService{}
	method := strings.TrimSpace(ctx.URLParam("method"))
	handlerName := strings.TrimSpace(ctx.URLParam("handlerName"))
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := actionService.Paging(offset, limit,handlerName,method)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "01030200", result)
}

func (self *ActionController)Create(ctx *iris.Context) {
	actionService := mngService.ActionService{}
	action := mngModel.Action{}
	if err := ctx.ReadJSON(&action); err != nil {
		common.Render(ctx, "01030301", err)
		return
	}
	if err := actionService.Create(&action); err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "01030300", action)
}

func (self *ActionController)Delete(ctx *iris.Context) {
	actionService := mngService.ActionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := actionService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "01030400", nil)
}

func (self *ActionController)Update(ctx *iris.Context) {
	actionService := mngService.ActionService{}
	action := mngModel.Action{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = actionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	err = ctx.ReadJSON(&action)
	if err != nil {
		common.Render(ctx, "01030501", err)
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
		return
	}
	common.Render(ctx, "01030500", result)
}
