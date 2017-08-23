package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type RoleController struct {

}

func (self *RoleController)GetAll(ctx *iris.Context) {
	roleService := permission.RoleService{}
	roleList, err := roleService.GetAll()
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27050200", roleList)
	return
}
func (self *RoleController)Create(ctx *iris.Context) {
	roleService := permission.RoleService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27050301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27050302", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "27050303", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	status := params.Get("status").MustInt()
	role := model.Role{
		Name:name,
		Description:description,
		Status:status,
	}
	entity, err := roleService.Create(&role)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27050300", entity)
}
func (self *RoleController)Delete(ctx *iris.Context) {
	roleService := permission.RoleService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := roleService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27050400", nil)
}
func (self *RoleController)GetByID(ctx *iris.Context) {
	roleService := permission.RoleService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	role, err := roleService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27050100", role)
}
func (self *RoleController)Update(ctx *iris.Context) {
	roleService := permission.RoleService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27050501", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	role, err := roleService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27050502", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "27050503", nil)
		return
	}
	description, e := params.CheckGet("description")
	if e {
		role.Description = strings.TrimSpace(description.MustString())
	}
	status := params.Get("status").MustInt(0)
	role.Name = name
	role.Status = status
	entity, err := roleService.Update(role)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27050500", entity)
}

func (self *RoleController)AssignPermissions(ctx *iris.Context) {
	roleService := permission.RoleService{}
	rolePermissionRelService := permission.RolePermissionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = roleService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	permissionIDs := make([]int, 0)
	if err := ctx.ReadJSON(&permissionIDs); err != nil {
		common.Render(ctx, "27050601", err)
		return
	}
	result, err := rolePermissionRelService.AssignPermissions(id, permissionIDs)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27050600", result)
}

func (self *RoleController)GetPermissions(ctx *iris.Context) {
	roleService := permission.RoleService{}
	rolePermissionRelService := permission.RolePermissionRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = roleService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	result, err := rolePermissionRelService.GetPermissionIDsByRoleIDs(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27050700", result)
}
