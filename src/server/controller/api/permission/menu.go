package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type MenuController struct {

}

func (self *MenuController)GetByID(ctx *iris.Context) {
	menuService := permission.MenuService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "01040101", err)
		return
	}
	menu, err := menuService.GetByID(id)
	if err != nil {
		common.Render(ctx, "01040102", err)
	}
	common.Render(ctx, "01040100", menu)
}

func (self *MenuController)Paging(ctx *iris.Context) {
	menuService := permission.MenuService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := menuService.Paging(offset, limit)
	if err != nil {
		common.Render(ctx, "01040201", err)
		return
	}
	common.Render(ctx, "01040200", result)
	return
}

func (self *MenuController)Create(ctx *iris.Context) {
	menuService := permission.MenuService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "01040301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	level := params.Get("level").MustInt()
	url := strings.TrimSpace(params.Get("url").MustString())
	parentID := params.Get("parentId").MustInt()
	icon := params.Get("icon").MustString()
	status := params.Get("status").MustInt()
	menu := model.Menu{
		Name:name,
		Level:level,
		ParentID:parentID,
		Url:url,
		Status:status,
		Icon:icon,
	}
	entity, err := menuService.Create(&menu)
	if err != nil {
		common.Render(ctx, "01040302", err)
		return
	}
	common.Render(ctx, "01040300", entity)
}

func (self *MenuController)Update(ctx *iris.Context) {
	menuService := permission.MenuService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "01040501", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "01040502", err)
		return
	}

	menu, err := menuService.GetByID(id)
	if err != nil {
		common.Render(ctx, "01040503", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == ""{
		common.Render(ctx, "01040504", nil)
		return
	}
	level := params.Get("level").MustInt(1)
	url, e := params.CheckGet("url")
	if e {
		menu.Url = strings.TrimSpace(url.MustString())
	}
	parentID := params.Get("parentId").MustInt(0)
	icon := strings.TrimSpace(params.Get("icon").MustString())
	if icon == ""{
		common.Render(ctx, "01040505", nil)
		return
	}
	status := params.Get("status").MustInt(0)
	menu.Name = name
	menu.Level = level
	menu.ParentID = parentID
	menu.Status = status
	menu.Icon = icon
	entity, err := menuService.Update(menu)
	if err != nil {
		common.Render(ctx, "01040506", err)
		return
	}
	common.Render(ctx, "01040500", entity)
}

func (self *MenuController)Delete(ctx *iris.Context) {
	menuService := permission.MenuService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "01040401", err)
		return
	}
	if err := menuService.Delete(id); err != nil {
		common.Render(ctx, "01040402", err)
	}
	common.Render(ctx, "01040400", nil)
}

func (self *MenuController)RearrangePosition(ctx *iris.Context) {
	menuService := permission.MenuService{}
	menuList := make([]*model.Menu, 0)
	//params := simplejson.New()
	if err := ctx.ReadJSON(&menuList); err != nil {
		common.Render(ctx, "01040601", err)
		return
	}

	entity, err := menuService.RearrangePosition(&menuList)
	if err != nil {
		common.Render(ctx, "01040601", err)
		return
	}
	common.Render(ctx, "01040600", entity)
}
