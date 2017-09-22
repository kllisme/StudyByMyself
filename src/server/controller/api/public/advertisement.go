package soda_manager

import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	model "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"time"
	"path/filepath"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/kit/util"
)

type AdvertisementController struct {

}

func (self *AdvertisementController)GetByID(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	adPositionService := mngService.ADPositionService{}
	applicationService := mngService.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04040101", err)
		return
	}
	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04040102", err)
		return
	}

	adPosition, err := adPositionService.GetByID(advertisement.AdPositionID)
	if err != nil {
		common.Render(ctx, "04040103", err)
		return
	}
	advertisement.AdPositionName = adPosition.Name

	app, err := applicationService.GetByID(adPosition.APPID)
	if err != nil {
		common.Render(ctx, "04040104", err)
		return
	}
	advertisement.APPName = app.Name
	advertisement.APPID = app.ID

	common.Render(ctx, "04040100", advertisement)
}

func (self *AdvertisementController)Paging(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	adPositionService := mngService.ADPositionService{}
	applicationService := mngService.ApplicationService{}
	var adPositionIDs []int
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	name := strings.TrimSpace(ctx.URLParam("name"))
	appID, _ := ctx.URLParamInt("appId")
	adPositionID, _ := ctx.URLParamInt("adPositionId")
	if adPositionID != 0 {
		adPositionIDs = []int{adPositionID}
	} else if appID != 0 {
		list, err := adPositionService.GetIDsByAPPID(appID)
		if err != nil {
			common.Render(ctx, "04040203", err)
			return
		}
		if len(list) == 0 {
			adPositionIDs = []int{0}
		} else {
			adPositionIDs = list
		}
	} else {
		adPositionIDs = []int{}
	}
	start := strings.TrimSpace(ctx.URLParam("startAt"))
	end := strings.TrimSpace(ctx.URLParam("endAt"))
	display, _ := ctx.URLParamInt("display")
	status, _ := ctx.URLParamInt("status")
	common.Logger.Debug(name)
	result, err := advertisementService.Paging("", name, adPositionIDs, start, end, display, status, offset, limit)
	if err != nil {
		common.Render(ctx, "04040201", err)
		return
	}

	adList := result.Objects.([]*model.Advertisement)
	for _, ad := range adList {
		adPosition, err := adPositionService.GetByID(ad.AdPositionID)
		if err != nil {
			common.Render(ctx, "04040204", err)
			return
		}
		ad.AdPositionName = adPosition.Name

		app, err := applicationService.GetByID(adPosition.APPID)
		if err != nil {
			common.Render(ctx, "04040205", err)
			return
		}
		ad.APPName = app.Name
		ad.APPID = app.ID
	}

	common.Render(ctx, "04040200", result)
	return
}

func (self *AdvertisementController)Create(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	adPositionService := mngService.ADPositionService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04040301", err)
		return
	}

	adPositionID := params.Get("adPositionId").MustInt()
	if adPositionID == 0 {
		common.Render(ctx, "04040302", nil)
		return
	}
	adPosition, err := adPositionService.GetByID(adPositionID)
	if err != nil {
		common.Render(ctx, "04040303", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04040304", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "04040305", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "04040306", nil)
		return
	} else if functions.CountRune(title) > 20 {
		common.Render(ctx, "04040307", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "04040308", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "04040309", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04040310", err)
		return
	}

	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04040311", err)
		return
	}

	display := params.Get("displayStrategy").MustInt()
	displayParams := ""
	if display == 2 {
		if adPosition.IdentifyNeeded == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
			if displayParams == "" {
				common.Render(ctx, "04040312", err)
				return
			}
		} else {
			common.Render(ctx, "04040313", err)
			return
		}
	}
	status := params.Get("status").MustInt()

	advertisement := model.Advertisement{
		Name:name,
		DisplayParams:displayParams,
		Status:status,
		AdPositionID:adPositionID,
		Title:title,
		URL:url,
		Image:image,
		StartedAt:startAt,
		EndedAt:endAt,
		DisplayStrategy:display,
	}

	if p, err := advertisementService.Paging(name, "", []int{}, "", "", 0, 0, 0, 0); err != nil {
		common.Render(ctx, "04040315", err)
		return
	} else {
		if p.Pagination.Total != 0 {
			common.Render(ctx, "04040316", nil)
			return
		}
	}

	entity, err := advertisementService.Create(&advertisement)
	if err != nil {
		common.Render(ctx, "04040314", err)
		return
	}
	common.Render(ctx, "04040300", entity)
}

func (self *AdvertisementController)Update(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	adPositionService := mngService.ADPositionService{}
	params := simplejson.New()

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04040501", err)
		return
	}

	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04040502", err)
		return
	}

	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04040503", err)
		return
	}


	adPosition, err := adPositionService.GetByID(advertisement.AdPositionID)
	if err != nil {
		common.Render(ctx, "04040505", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04040506", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "04040507", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "04040508", nil)
		return
	} else if functions.CountRune(title) > 20 {
		common.Render(ctx, "04040509", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "04040510", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "04040511", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04040512", err)
		return
	}
	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04040513", err)
		return
	}

	display := params.Get("displayStrategy").MustInt()
	displayParams := ""
	if display == 2 {
		if adPosition.IdentifyNeeded == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
			if displayParams == "" {
				common.Render(ctx, "04040514", err)
				return
			}
		} else {
			common.Render(ctx, "04040515", err)
			return
		}
	}
	status := params.Get("status").MustInt()

	if p, err := advertisementService.Paging(name, "", []int{}, "", "", 0, 0, 0, 0); err != nil {
		common.Render(ctx, "04040517", err)
		return
	} else {
		if p.Pagination.Total != 0 && name != advertisement.Name {
			common.Render(ctx, "04040518", nil)
			return
		}
	}

	advertisement.URL = url
	advertisement.Image = image
	advertisement.DisplayStrategy = display
	advertisement.DisplayParams = displayParams
	advertisement.Status = status
	advertisement.StartedAt = startAt
	advertisement.EndedAt = endAt
	advertisement.Title = title
	advertisement.Name = name
	entity, err := advertisementService.Update(advertisement)
	if err != nil {
		common.Render(ctx, "04040516", err)
		return
	}
	common.Render(ctx, "04040500", entity)
}

func (self *AdvertisementController)Delete(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04040401", err)
		return
	}
	if err := advertisementService.Delete(id); err != nil {
		common.Render(ctx, "04040402", err)
		return
	}
	common.Render(ctx, "04040400", nil)
}

func (self *AdvertisementController)BatchUpdateOrder(ctx *iris.Context) {
	advertisementService := mngService.AdvertisementService{}
	adList := make([]*model.Advertisement, 0)
	//params := simplejson.New()
	if err := ctx.ReadJSON(&adList); err != nil {
		common.Render(ctx, "04040601", err)
		return
	}

	entity, err := advertisementService.BatchUpdateOrder(&adList)
	if err != nil {
		common.Render(ctx, "04040602", err)
		return
	}
	common.Render(ctx, "04040600", entity)
}

func (self *AdvertisementController)SaveImage(ctx *iris.Context) {
	ctx.Set("isUpload", true)
	object := viper.GetString("resource.oss.object.ad")
	formFile, err := ctx.FormFile("file")
	fileName := formFile.Filename
	fileExt := filepath.Ext(fileName)
	if fileExt != ".jpg" && fileExt != ".png" {
		common.Render(ctx, "04040701", nil)
		return
	}

	shortPath, err := util.Upload(formFile, object)
	if err != nil {
		common.Render(ctx, "04040702", err)
		return
	}
	// domain := viper.GetString("resource.oss.domain")
	common.Render(ctx, "04040700", shortPath)
}
