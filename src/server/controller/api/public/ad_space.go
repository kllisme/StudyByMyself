package public

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	model "maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type ADSpaceController struct {

}

func (self *ADSpaceController)GetByID(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	adSpace, err := adSpaceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	app, err := applicationService.GetByID(adSpace.APPID)
	if err != nil {
		common.Render(ctx, "000002", nil)
	}
	adSpace.APPName = app.Name
	common.Render(ctx, "27070100", adSpace)
}

func (self *ADSpaceController)Paging(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}

	appID, _ := ctx.URLParamInt("app_id")
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := adSpaceService.Paging(appID, page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	adSpaceList := result.Objects.([]*model.ADSpace)
	for _, value := range adSpaceList {
		app, err := applicationService.GetByID(value.APPID)
		if err != nil {
			common.Render(ctx, "000002", nil)
			return
		}
		value.APPName = app.Name
	}
	common.Render(ctx, "27070200", result)
}

func (self *ADSpaceController)Create(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070301", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	}

	description := strings.TrimSpace(params.Get("description").MustString())
	appID := params.Get("appId").MustInt()
	if appID == 0 {
		common.Render(ctx, "27070303", nil)
		return
	}
	identifyNeeded := params.Get("identifyNeeded").MustInt()
	adSpace := model.ADSpace{
		Name:name,
		APPID:appID,
		Description:description,
		IdentifyNeeded:identifyNeeded,
	}
	entity, err := adSpaceService.Create(&adSpace)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *ADSpaceController)Update(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
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

	adSpace, err := adSpaceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	appID := params.Get("appId").MustInt()
	if appID == 0 {
		common.Render(ctx, "27070303", nil)
		return
	}
	identifyNeeded := params.Get("identifyNeeded").MustInt()
	adSpace.Name = name
	adSpace.Description = description
	adSpace.IdentifyNeeded = identifyNeeded
	adSpace.APPID = appID
	entity, err := adSpaceService.Update(adSpace)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *ADSpaceController)Delete(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	AdvertisementService := public.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	//TODO 测试如果没有获取到匹配的 adlist 时的表现
	adList,err := AdvertisementService.GetListByLocationID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if len(*adList) !=0 {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := adSpaceService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27070400", nil)
}
