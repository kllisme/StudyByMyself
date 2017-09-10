package permission

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	model "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type ElementController struct {

}

func (self *ElementController)GetByID(ctx *iris.Context) {
	elementService := permission.ElementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	element, err := elementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070100", element)
}

func (self *ElementController)Paging(ctx *iris.Context) {
	elementService := permission.ElementService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	name := ctx.URLParam("name")
	result, err := elementService.Paging(name, offset, limit)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070200", result)
	return
}

func (self *ElementController)Create(ctx *iris.Context) {
	elementService := permission.ElementService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	}
	reference := strings.TrimSpace(params.Get("reference").MustString())
	//if reference == "" {
	//	common.Render(ctx, "27070303", nil)
	//	return
	//}
	element := model.Element{
		Name:name,
		Reference:reference,
	}
	entity, err := elementService.Create(&element)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *ElementController)Update(ctx *iris.Context) {
	elementService := permission.ElementService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070501", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	element, err := elementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}
	reference := strings.TrimSpace(params.Get("reference").MustString())
	//if reference == "" {
	//	common.Render(ctx, "27070503", nil)
	//	return
	//}
	element.Name = name
	element.Reference = reference
	entity, err := elementService.Update(element)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *ElementController)Delete(ctx *iris.Context) {
	elementService := permission.ElementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := elementService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070400", nil)
}
