package permission

import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type PermissionController struct {

}

func (self *PermissionController)GetByID(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_permission, err := permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27060100", _permission)
}

func (self *PermissionController)Paging(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	categoryID, _ := ctx.URLParamInt("categoryId")
	result, err := permissionService.Paging(categoryID, offset, limit)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060200", result)
	return
}

func (self *PermissionController)Create(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27060301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27060302", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "27060303", nil)
		return
	}
	categoryID := params.Get("categoryId").MustInt()
	status := params.Get("status").MustInt()
	_permission := mngModel.Permission{
		Name:name,
		CategoryID:categoryID,
		Status:status,
	}
	entity, err := permissionService.Create(&_permission)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060300", entity)
}

func (self *PermissionController)Update(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permission := mngModel.Permission{}

	if err := ctx.ReadJSON(&permission); err != nil {
		common.Render(ctx, "27060501", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	permission.Name = strings.TrimSpace(permission.Name)
	if permission.Name == "" {
		common.Render(ctx, "27060502", nil)
		return
	} else if functions.CountRune(permission.Name) > 20 {
		common.Render(ctx, "27060503", nil)
		return
	}
	permission.ID = id
	entity, err := permissionService.Update(&permission)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060500", entity)
}

func (self *PermissionController)Delete(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := permissionService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27060400", nil)
}

func (self *PermissionController)AssignMenus(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionMenuRelService := mngService.PermissionMenuRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	menuIDs := make([]int, 0)
	if err := ctx.ReadJSON(&menuIDs); err != nil {
		common.Render(ctx, "27060601", err)
		return
	}
	result, err := permissionMenuRelService.AssignMenus(id, menuIDs)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060600", result)
}

func (self *PermissionController)GetMenus(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionMenuRelService := mngService.PermissionMenuRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	result, err := permissionMenuRelService.GetMenuIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060700", result)
}

func (self *PermissionController)AssignActions(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionActionRelService := mngService.PermissionActionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	actionIDs := make([]int, 0)
	if err := ctx.ReadJSON(&actionIDs); err != nil {
		common.Render(ctx, "27060801", err)
		return
	}
	result, err := permissionActionRelService.AssignActions(id, actionIDs)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060800", result)
}

func (self *PermissionController)GetActions(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionActionRelService := mngService.PermissionActionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	result, err := permissionActionRelService.GetActionIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27060900", result)
}

func (self *PermissionController)AssignElements(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionElementRelService := mngService.PermissionElementRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	elementIDs := make([]int, 0)
	if err := ctx.ReadJSON(&elementIDs); err != nil {
		common.Render(ctx, "27061001", err)
		return
	}
	result, err := permissionElementRelService.AssignElements(id, elementIDs)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27061000", result)
}

func (self *PermissionController)GetElements(ctx *iris.Context) {
	permissionService := mngService.PermissionService{}
	permissionElementRelService := mngService.PermissionElementRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = permissionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	result, err := permissionElementRelService.GetElementIDsByPermissionIDs(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27061100", result)
}
