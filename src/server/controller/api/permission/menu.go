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
		common.Render(ctx, "000003", nil)
		return
	}
	menu, err := menuService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "", menu)
}

func (self *MenuController)Paging(ctx *iris.Context) {
	menuService := permission.MenuService{}
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := menuService.Paging(page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020300", result)
	return
}

func (self *MenuController)Create(ctx *iris.Context) {
	menuService := permission.MenuService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020201", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	//name := strings.TrimSpace(params.MustString("name"))
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
		common.Render(ctx, "27010103", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *MenuController)Update(ctx *iris.Context) {
	menuService := permission.MenuService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020201", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	//name := strings.TrimSpace(params.MustString("name"))
	level := params.Get("level").MustInt(1)
	url := strings.TrimSpace(params.Get("url").MustString())
	parentID := params.Get("parentId").MustInt(0)
	icon := params.Get("icon").MustString()
	status := params.Get("status").MustInt(0)
	menu := model.Menu{
		Name:name,
		Level:level,
		ParentID:parentID,
		Url:url,
		Status:status,
		Icon:icon,
	}
	menu.ID = id
	entity, err := menuService.Update(&menu)
	if err != nil {
		common.Render(ctx, "27010103", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *MenuController)Delete(ctx *iris.Context) {
	menuService := permission.MenuService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := menuService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "", nil)
}
