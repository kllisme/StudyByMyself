package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type PermissionController struct {

}

func (self *PermissionController)GetByID(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_permission, err := permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27060100", _permission)
}

func (self *PermissionController)Paging(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	_type, _ := ctx.URLParamInt("type")
	result, err := permissionService.Paging(_type, page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060200", result)
	return
}

func (self *PermissionController)Create(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27060301", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27060302", nil)
		return
	}
	_type := params.Get("type").MustInt()
	status := params.Get("status").MustInt()
	_permission := model.Permission{
		Name:name,
		Type:_type,
		Status:status,
	}
	entity, err := permissionService.Create(&_permission)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060300", entity)
}

func (self *PermissionController)Update(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27060501", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	_permission, err := permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27060502", nil)
		return
	}
	_type := params.Get("type").MustInt(0)
	status := params.Get("status").MustInt(0)
	_permission.Name = name
	_permission.Status = status
	_permission.Type = _type
	entity, err := permissionService.Update(_permission)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060500", entity)
}

func (self *PermissionController)Delete(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := permissionService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27060400", nil)
}

func (self *PermissionController)AssignElements(ctx *iris.Context) {

}

func (self *PermissionController)AssignMenus(ctx *iris.Context) {

}

func (self *PermissionController)AssignActions(ctx *iris.Context) {

}
