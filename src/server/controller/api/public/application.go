package soda_manager
import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type ApplicationController struct {

}

func (self *ApplicationController)GetByID(ctx *iris.Context) {
	applicationService := mngService.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020101", err)
		return
	}
	application, err := applicationService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020102", err)
		return
	}
	common.Render(ctx, "04020100", application)
}

func (self *ApplicationController)Paging(ctx *iris.Context) {
	applicationService := mngService.ApplicationService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := applicationService.Paging(offset, limit)
	if err != nil {
		common.Render(ctx, "04020201", nil)
		return
	}
	common.Render(ctx, "04020200", result)
	return
}

func (self *ApplicationController)Create(ctx *iris.Context) {
	applicationService := mngService.ApplicationService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04020301", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04020302", nil)
		return
	}else if functions.CountRune(name) > 10 {
		common.Render(ctx, "04020303", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "04020304", nil)
		return
	}
	application := mngModel.Application{
		Name:name,
		Description:description,
	}
	entity, err := applicationService.Create(&application)
	if err != nil {
		common.Render(ctx, "04020305", nil)
		return
	}
	common.Render(ctx, "04020300", entity)
}

func (self *ApplicationController)Update(ctx *iris.Context) {
	applicationService := mngService.ApplicationService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04020507", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020501", nil)
		return
	}

	application, err := applicationService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020502", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04020503", nil)
		return
	}else if functions.CountRune(name) > 10 {
		common.Render(ctx, "04020504", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "04020505", nil)
		return
	}
	application.Name = name
	application.Description = description
	entity, err := applicationService.Update(application)
	if err != nil {
		common.Render(ctx, "04020506", nil)
		return
	}
	common.Render(ctx, "04020500", entity)
}

func (self *ApplicationController)Delete(ctx *iris.Context) {
	applicationService := mngService.ApplicationService{}
	adPositionService := mngService.ADPositionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020401", nil)
		return
	}

	adList, err := adPositionService.GetIDsByAPPID(id)
	if err != nil {
		common.Render(ctx, "04020402", err)
		return
	}
	if len(adList) != 0 {
		common.Render(ctx, "04020403", nil)
		return
	}

	if err := applicationService.Delete(id); err != nil {
		common.Render(ctx, "04020404", nil)
		return
	}
	common.Render(ctx, "04020400", nil)
}
