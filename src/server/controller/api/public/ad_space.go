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

type ADSpaceController struct {

}

func (self *ADSpaceController)GetByID(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030101", err)
		return
	}
	adSpace, err := adSpaceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04030102", err)
	}
	app, err := applicationService.GetByID(adSpace.APPID)
	if err != nil {
		common.Render(ctx, "04030103", err)
	}
	adSpace.APPName = app.Name
	common.Render(ctx, "04030100", adSpace)
}

func (self *ADSpaceController)Paging(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}

	appID, _ := ctx.URLParamInt("appId")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := adSpaceService.Paging("", appID, offset, limit)
	if err != nil {
		common.Render(ctx, "04030201", err)
		return
	}
	adSpaceList := result.Objects.([]*model.ADSpace)
	for _, value := range adSpaceList {
		app, err := applicationService.GetByID(value.APPID)
		if err != nil {
			common.Render(ctx, "04030202", err)
			return
		}
		value.APPName = app.Name
	}
	common.Render(ctx, "04030200", result)
}

func (self *ADSpaceController)Create(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04030301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04030302", nil)
		return
	} else if functions.CountRune(name) > 10 {
		common.Render(ctx, "04030303", nil)
		return
	}

	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "04030304", nil)
		return
	}
	appID := params.Get("appId").MustInt()
	if appID == 0 {
		common.Render(ctx, "04030305", nil)
		return
	}
	identifyNeeded := params.Get("identifyNeeded").MustInt()
	standard := strings.TrimSpace(params.Get("standard").MustString())
	if standard == "" {
		common.Render(ctx, "04030307", nil)
		return
	}
	adSpace := model.ADSpace{
		Name:name,
		APPID:appID,
		Description:description,
		IdentifyNeeded:identifyNeeded,
		Standard:standard,
	}

	if p, err := adSpaceService.Paging(name, appID, 0, 0); err != nil {
		common.Render(ctx, "04030308", err)
		return
	} else {
		if p.Pagination.Total != 0 {
			common.Render(ctx, "04030309", nil)
			return
		}
	}

	entity, err := adSpaceService.Create(&adSpace)
	if err != nil {
		common.Render(ctx, "04030306", err)
		return
	}
	common.Render(ctx, "04030300", entity)
}

func (self *ADSpaceController)Update(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030501", err)
		return
	}

	adSpace, err := adSpaceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04030502", err)
		return
	}

	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04030503", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04030504", nil)
		return
	} else if functions.CountRune(name) > 10 {
		common.Render(ctx, "04030505", nil)
		return
	}
	description := strings.TrimSpace(params.Get("description").MustString())
	if functions.CountRune(description) > 50 {
		common.Render(ctx, "04030506", nil)
		return
	}
	//appID := params.Get("appId").MustInt()
	//if appID == 0 {
	//	common.Render(ctx, "04030507", nil)
	//	return
	//}
	standard := strings.TrimSpace(params.Get("standard").MustString())
	if standard == "" {
		common.Render(ctx, "04030509", nil)
		return
	}
	identifyNeeded := params.Get("identifyNeeded").MustInt()

	if p, err := adSpaceService.Paging(name, adSpace.APPID, 0, 0); err != nil {
		common.Render(ctx, "04030510", err)
		return
	} else {
		if p.Pagination.Total != 0 && name != adSpace.Name {
			common.Render(ctx, "04030511", nil)
			return
		}
	}

	adSpace.Name = name
	adSpace.Description = description
	adSpace.IdentifyNeeded = identifyNeeded
	//adSpace.APPID = appID
	adSpace.Standard = standard
		entity, err := adSpaceService.Update(adSpace)
	if err != nil {
		common.Render(ctx, "04030508", err)
		return
	}
	common.Render(ctx, "04030500", entity)
}

func (self *ADSpaceController)Delete(ctx *iris.Context) {
	adSpaceService := public.ADSpaceService{}
	AdvertisementService := public.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030401", err)
		return
	}
	adList, err := AdvertisementService.GetListByLocationID(id)
	if err != nil {
		common.Render(ctx, "04030402", err)
		return
	}
	if len(*adList) != 0 {
		common.Render(ctx, "04030403", nil)
		return
	}
	if err := adSpaceService.Delete(id); err != nil {
		common.Render(ctx, "04030404", err)
	}
	common.Render(ctx, "04030400", nil)
}
