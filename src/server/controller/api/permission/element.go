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
		common.Render(ctx, "000003", nil)
		return
	}
	element, err := elementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040100", element)
}

func (self *ElementController)Paging(ctx *iris.Context) {
	elementService := permission.ElementService{}
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := elementService.Paging(page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27040200", result)
	return
}

func (self *ElementController)Create(ctx *iris.Context) {
	elementService := permission.ElementService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27040301", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "", nil)
		return
	}
	reference := strings.TrimSpace(params.Get("reference").MustString())
	if reference == "" {
		common.Render(ctx, "", nil)
		return
	}
	element := model.Element{
		Name:name,
		Reference:reference,
	}
	entity, err := elementService.Create(&element)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27040300", entity)
}

func (self *ElementController)Update(ctx *iris.Context) {
	elementService := permission.ElementService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27040501", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	element, err := elementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == ""{
		common.Render(ctx, "27020502", nil)
		return
	}
	reference := strings.TrimSpace(params.Get("reference").MustString())
	if reference == "" {
		common.Render(ctx, "", nil)
		return
	}
	element.Name = name
	element.Reference = reference
	entity, err := elementService.Update(element)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27040500", entity)
}

func (self *ElementController)Delete(ctx *iris.Context) {
	elementService := permission.ElementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := elementService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27040400", nil)
}
