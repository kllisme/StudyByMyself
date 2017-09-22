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

type ADPositionController struct {

}

func (self *ADPositionController)GetByID(ctx *iris.Context) {
	adPositionService := mngService.ADPositionService{}
	applicationService := mngService.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030101", err)
		return
	}
	adPosition, err := adPositionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04030102", err)
		return
	}
	app, err := applicationService.GetByID(adPosition.APPID)
	if err != nil {
		common.Render(ctx, "04030103", err)
		return
	}
	adPosition.APPName = app.Name
	common.Render(ctx, "04030100", adPosition)
}

func (self *ADPositionController)Paging(ctx *iris.Context) {
	adPositionService := mngService.ADPositionService{}
	applicationService := mngService.ApplicationService{}

	appID, _ := ctx.URLParamInt("appId")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := adPositionService.Paging("", appID, offset, limit)
	if err != nil {
		common.Render(ctx, "04030201", err)
		return
	}
	adPositionList := result.Objects.([]*mngModel.ADPosition)
	for _, value := range adPositionList {
		app, err := applicationService.GetByID(value.APPID)
		if err != nil {
			common.Render(ctx, "04030202", err)
			return
		}
		value.APPName = app.Name
	}
	common.Render(ctx, "04030200", result)
}

func (self *ADPositionController)Create(ctx *iris.Context) {
	adPositionService := mngService.ADPositionService{}
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
	adPosition := mngModel.ADPosition{
		Name:name,
		APPID:appID,
		Description:description,
		IdentifyNeeded:identifyNeeded,
		Standard:standard,
	}

	if p, err := adPositionService.Paging(name, appID, 0, 0); err != nil {
		common.Render(ctx, "04030308", err)
		return
	} else {
		if p.Pagination.Total != 0 {
			common.Render(ctx, "04030309", nil)
			return
		}
	}

	entity, err := adPositionService.Create(&adPosition)
	if err != nil {
		common.Render(ctx, "04030306", err)
		return
	}
	common.Render(ctx, "04030300", entity)
}

func (self *ADPositionController)Update(ctx *iris.Context) {
	adPositionService := mngService.ADPositionService{}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030501", err)
		return
	}

	adPosition, err := adPositionService.GetByID(id)
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

	if p, err := adPositionService.Paging(name, adPosition.APPID, 0, 0); err != nil {
		common.Render(ctx, "04030510", err)
		return
	} else {
		if p.Pagination.Total != 0 && name != adPosition.Name {
			common.Render(ctx, "04030511", nil)
			return
		}
	}

	adPosition.Name = name
	adPosition.Description = description
	adPosition.IdentifyNeeded = identifyNeeded
	//adPosition.APPID = appID
	adPosition.Standard = standard
		entity, err := adPositionService.Update(adPosition)
	if err != nil {
		common.Render(ctx, "04030508", err)
		return
	}
	common.Render(ctx, "04030500", entity)
}

func (self *ADPositionController)Delete(ctx *iris.Context) {
	adPositionService := mngService.ADPositionService{}
	AdvertisementService := mngService.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04030401", err)
		return
	}
	adList, err := AdvertisementService.GetListByADPositionID(id)
	if err != nil {
		common.Render(ctx, "04030402", err)
		return
	}
	if len(*adList) != 0 {
		common.Render(ctx, "04030403", nil)
		return
	}
	if err := adPositionService.Delete(id); err != nil {
		common.Render(ctx, "04030404", err)
		return
	}
	common.Render(ctx, "04030400", nil)
}
