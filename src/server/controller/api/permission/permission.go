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
	categoryID, _ := ctx.URLParamInt("category_id")
	result, err := permissionService.Paging(categoryID, page, perPage)
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
	categoryID := params.Get("categoryId").MustInt()
	status := params.Get("status").MustInt()
	_permission := model.Permission{
		Name:name,
		CategoryID:categoryID,
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
	permission:= model.Permission{}

	if err := ctx.ReadJSON(&permission); err != nil {
		common.Render(ctx, "27060501", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	permission.Name = strings.TrimSpace(permission.Name)
	if permission.Name == "" {
		common.Render(ctx, "27060502", nil)
		return
	}
	permission.ID = id
	entity, err := permissionService.Update(&permission)
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

func (self *PermissionController)AssignMenus(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionMenuRelService := permission.PermissionMenuRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	menuIDs := make([]int, 0)
	if err := ctx.ReadJSON(&menuIDs); err != nil {
		common.Render(ctx, "27060601", nil)
		return
	}
	result, err := permissionMenuRelService.AssignMenus(id, menuIDs)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060600", result)
}

func (self *PermissionController)GetMenus(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionMenuRelService := permission.PermissionMenuRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	result, err := permissionMenuRelService.GetMenuIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060700", result)
}

func (self *PermissionController)AssignActions(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionActionRelService := permission.PermissionActionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	actionIDs := make([]int, 0)
	if err := ctx.ReadJSON(&actionIDs); err != nil {
		common.Render(ctx, "27060801", nil)
		return
	}
	result, err := permissionActionRelService.AssignActions(id, actionIDs)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060800", result)
}

func (self *PermissionController)GetActions(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionActionRelService := permission.PermissionActionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	result, err := permissionActionRelService.GetActionIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27060900", result)
}

func (self *PermissionController)AssignElements(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionElementRelService := permission.PermissionElementRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	elementIDs := make([]int, 0)
	if err := ctx.ReadJSON(&elementIDs); err != nil {
		common.Render(ctx, "27061001", nil)
		return
	}
	result, err := permissionElementRelService.AssignElements(id, elementIDs)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27061000", result)
}

func (self *PermissionController)GetElements(ctx *iris.Context) {
	permissionService := permission.PermissionService{}
	permissionElementRelService := permission.PermissionElementRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	result, err := permissionElementRelService.GetElementIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27061100", result)
}
