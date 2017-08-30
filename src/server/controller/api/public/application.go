package public
import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	model "maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type ApplicationController struct {

}

func (self *ApplicationController)GetByID(ctx *iris.Context) {
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	application, err := applicationService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070100", application)
}

func (self *ApplicationController)Paging(ctx *iris.Context) {
	applicationService := public.ApplicationService{}
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := applicationService.Paging(page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070200", result)
	return
}

func (self *ApplicationController)Create(ctx *iris.Context) {
	applicationService := public.ApplicationService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070301", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	}else if functions.CountRune(name) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "27070303", nil)
		return
	}
	application := model.Application{
		Name:name,
		Description:description,
	}
	entity, err := applicationService.Create(&application)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *ApplicationController)Update(ctx *iris.Context) {
	applicationService := public.ApplicationService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070501", nil)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	application, err := applicationService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	}else if functions.CountRune(name) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "27070303", nil)
		return
	}
	application.Name = name
	application.Description = description
	entity, err := applicationService.Update(application)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *ApplicationController)Delete(ctx *iris.Context) {
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := applicationService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27070400", nil)
}
