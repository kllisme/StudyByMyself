package public

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	model "maizuo.com/soda/erp/api/src/server/model/public"
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
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020101", err)
		return
	}
	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020102", err)
		return
	}

	adSpace, err := adSpaceService.GetByID(advertisement.LocationID)
	if err != nil {
		common.Render(ctx, "04020103", err)
		return
	}
	advertisement.LocationName = adSpace.Name

	app, err := applicationService.GetByID(adSpace.APPID)
	if err != nil {
		common.Render(ctx, "04020104", err)
		return
	}
	advertisement.APPName = app.Name
	advertisement.APPID = app.ID

	common.Render(ctx, "04020100", advertisement)
}

func (self *AdvertisementController)Paging(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	var locationIDs []int
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	title := strings.TrimSpace(ctx.URLParam("title"))
	appID, _ := ctx.URLParamInt("app_id")
	locationID, _ := ctx.URLParamInt("location_id")
	if locationID != 0 {
		locationIDs = []int{locationID}
	} else if appID != 0 {
		list, err := adSpaceService.GetLocationIDs(appID)
		if err != nil {
			common.Render(ctx, "04020203", err)
			return
		}
		if len(list) == 0 {
			locationIDs = []int{0}
		} else {
			locationIDs = list
		}
	} else {
		locationIDs = []int{}
	}
	start := strings.TrimSpace(ctx.URLParam("started_at"))
	end := strings.TrimSpace(ctx.URLParam("ended_at"))
	display, _ := ctx.URLParamInt("display")
	status, _ := ctx.URLParamInt("status")
	result, err := advertisementService.Paging(title, locationIDs, start, end, display, status, page, perPage)
	if err != nil {
		common.Render(ctx, "04020201", err)
		return
	}

	adList := result.Objects.([]*model.Advertisement)
	for _, ad := range adList {
		adSpace, err := adSpaceService.GetByID(ad.LocationID)
		if err != nil {
			common.Render(ctx, "04020204", err)
			return
		}
		ad.LocationName = adSpace.Name

		app, err := applicationService.GetByID(adSpace.APPID)
		if err != nil {
			common.Render(ctx, "04020205", err)
			return
		}
		ad.APPName = app.Name
		ad.APPID = app.ID
	}

	common.Render(ctx, "04020200", result)
	return
}

func (self *AdvertisementController)Create(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04020301", err)
		return
	}

	locationID := params.Get("locationId").MustInt()
	if locationID == 0 {
		common.Render(ctx, "04020302", nil)
		return
	}
	adSpace, err := adSpaceService.GetByID(locationID)
	if err != nil {
		common.Render(ctx, "04020303", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04020304", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "04020305", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "04020306", nil)
		return
	} else if functions.CountRune(title) > 20 {
		common.Render(ctx, "04020307", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "04020308", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "04020309", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04020310", err)
		return
	}

	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04020311", err)
		return
	}

	display := params.Get("displayStrategy").MustInt()
	displayParams := ""
	if display == 2 {
		if adSpace.IdentifyNeeded == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
			if displayParams == "" {
				common.Render(ctx, "04020312", err)
				return
			}
		} else {
			common.Render(ctx, "04020313", err)
			return
		}
	}
	status := params.Get("status").MustInt()

	advertisement := model.Advertisement{
		Name:name,
		DisplayParams:displayParams,
		Status:status,
		LocationID:locationID,
		Title:title,
		URL:url,
		Image:image,
		StartedAt:startAt,
		EndedAt:endAt,
		DisplayStrategy:display,
	}
	entity, err := advertisementService.Create(&advertisement)
	if err != nil {
		common.Render(ctx, "04020314", err)
		return
	}
	common.Render(ctx, "04020300", entity)
}

func (self *AdvertisementController)Update(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	params := simplejson.New()

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020501", err)
		return
	}

	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020502", err)
		return
	}

	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04020503", err)
		return
	}

	locationID := params.Get("locationId").MustInt()
	if locationID == 0 {
		common.Render(ctx, "04020504", nil)
		return
	}
	adSpace, err := adSpaceService.GetByID(locationID)
	if err != nil {
		common.Render(ctx, "04020505", err)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04020506", nil)
		return
	} else if functions.CountRune(name) > 20 {
		common.Render(ctx, "04020507", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "04020508", nil)
		return
	} else if functions.CountRune(title) > 20 {
		common.Render(ctx, "04020509", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "04020510", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "04020511", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04020512", err)
		return
	}
	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endedAt").MustString()))
	if err != nil {
		common.Render(ctx, "04020513", err)
		return
	}

	display := params.Get("displayStrategy").MustInt()
	displayParams := ""
	if display == 2 {
		if adSpace.IdentifyNeeded == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
			if displayParams == "" {
				common.Render(ctx, "04020514", err)
				return
			}
		} else {
			common.Render(ctx, "04020515", err)
			return
		}
	}
	status := params.Get("status").MustInt()

	advertisement.URL = url
	advertisement.Image = image
	advertisement.DisplayStrategy = display
	advertisement.DisplayParams = displayParams
	advertisement.Status = status
	advertisement.StartedAt = startAt
	advertisement.EndedAt = endAt
	advertisement.LocationID = locationID
	advertisement.Title = title
	advertisement.Name = name
	entity, err := advertisementService.Update(advertisement)
	if err != nil {
		common.Render(ctx, "04020516", err)
		return
	}
	common.Render(ctx, "04020500", entity)
}

func (self *AdvertisementController)Delete(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020401", err)
		return
	}
	if err := advertisementService.Delete(id); err != nil {
		common.Render(ctx, "04020402", err)
	}
	common.Render(ctx, "04020400", nil)
}

func (self *AdvertisementController)BatchUpdateOrder(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adList := make([]*model.Advertisement, 0)
	//params := simplejson.New()
	if err := ctx.ReadJSON(&adList); err != nil {
		common.Render(ctx, "04020601", err)
		return
	}

	entity, err := advertisementService.BatchUpdateOrder(&adList)
	if err != nil {
		common.Render(ctx, "04020602", err)
		return
	}
	common.Render(ctx, "04020600", entity)
}

func (self *AdvertisementController)SaveImage(ctx *iris.Context) {
	ctx.Set("isUpload", true)
	object := viper.GetString("resource.oss.object.ad")
	formFile, err := ctx.FormFile("file")
	fileName := formFile.Filename
	fileExt := filepath.Ext(fileName)
	if fileExt != ".jpg" && fileExt != ".png" {
		common.Render(ctx, "04020701", nil)
		return
	}

	shortPath, err := util.Upload(formFile, object)
	if err != nil {
		common.Render(ctx, "04020702", err)
		return
	}
	// domain := viper.GetString("resource.oss.domain")
	common.Render(ctx, "04020700", shortPath)
}
