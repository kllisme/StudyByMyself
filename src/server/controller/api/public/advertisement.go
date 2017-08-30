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
)

type AdvertisementController struct {

}

func (self *AdvertisementController)GetByID(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}

	adSpace, err := adSpaceService.GetByID(advertisement.LocationID)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	advertisement.LocationName = adSpace.Name

	app, err := applicationService.GetByID(adSpace.APPID)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	advertisement.APPName = app.Name

	common.Render(ctx, "27070100", advertisement)
}

func (self *AdvertisementController)Paging(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	applicationService := public.ApplicationService{}
	locationIDs := []int{}
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	title := strings.TrimSpace(ctx.URLParam("title"))
	appID, _ := ctx.URLParamInt("app_id")
	locationID, _ := ctx.URLParamInt("location_id")
	if locationID != 0 {
		locationIDs[0] = locationID
	} else if appID != 0 {
		list, err := adSpaceService.GetLocationIDs(appID)
		if err != nil {
			common.Render(ctx, "27070200", err)
			return
		}
		locationIDs = list
	}
	start := strings.TrimSpace(ctx.URLParam("start"))
	end := strings.TrimSpace(ctx.URLParam("end"))
	display, _ := ctx.URLParamInt("display")
	status, _ := ctx.URLParamInt("status")

	result, err := advertisementService.Paging(title, locationIDs, start, end, display, status, page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}

	adList := result.Objects.([]*model.Advertisement)
	for _, ad := range adList {
		adSpace, err := adSpaceService.GetByID(ad.LocationID)
		if err != nil {
			common.Render(ctx, "000002", nil)
			return
		}
		ad.LocationName = adSpace.Name

		app, err := applicationService.GetByID(adSpace.APPID)
		if err != nil {
			common.Render(ctx, "000002", nil)
			return
		}
		ad.APPName = app.Name
	}

	common.Render(ctx, "27070200", result)
	return
}

func (self *AdvertisementController)Create(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	adSpaceService := public.ADSpaceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070301", nil)
		return
	}

	locationID := params.Get("locationId").MustInt()
	if locationID == 0 {
		common.Render(ctx, "27070302", nil)
		return
	}
	adSpace, err := adSpaceService.GetByID(locationID)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	} else if functions.CountRune(name) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "27070302", nil)
		return
	} else if functions.CountRune(title) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "27070302", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "27070302", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startAt").MustString()))
	if err != nil {
		common.Render(ctx, "27070302", nil)
		return
	}

	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endAt").MustString()))
	if err != nil {
		common.Render(ctx, "27070302", nil)
		return
	}

	display := 0
	displayParams := ""

	if adSpace.IdentifyNeeded == 1 {
		display = params.Get("displayStrategy").MustInt()

		if display == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
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
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *AdvertisementController)Update(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
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

	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	locationID := params.Get("locationId").MustInt()
	if locationID == 0 {
		common.Render(ctx, "27070302", nil)
		return
	}
	adSpace, err := adSpaceService.GetByID(locationID)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}

	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	} else if functions.CountRune(name) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}

	title := strings.TrimSpace(params.Get("title").MustString())
	if title == "" {
		common.Render(ctx, "27070302", nil)
		return
	} else if functions.CountRune(title) > 10 {
		common.Render(ctx, "27070302", nil)
		return
	}

	url := strings.TrimSpace(params.Get("url").MustString())
	if url == "" {
		common.Render(ctx, "27070302", nil)
		return
	}

	image := strings.TrimSpace(params.Get("image").MustString())
	if image == "" {
		common.Render(ctx, "27070302", nil)
		return
	}

	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("startAt").MustString()))
	if err != nil {
		common.Render(ctx, "27070302", nil)
		return
	}

	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(params.Get("endAt").MustString()))
	if err != nil {
		common.Render(ctx, "27070302", nil)
		return
	}

	display := 0
	displayParams := ""

	if adSpace.IdentifyNeeded == 1 {
		display = params.Get("displayStrategy").MustInt()

		if display == 1 {
			displayParams = strings.TrimSpace(params.Get("displayParams").MustString())
		}
	}
	status := params.Get("status").MustInt()

	advertisement.URL = url
	advertisement.Image = image
	advertisement.DisplayStrategy = display
	advertisement.DisplayParams =displayParams
	advertisement.Status =status
	advertisement.StartedAt = startAt
	advertisement.EndedAt =endAt
	advertisement.LocationID =locationID
	advertisement.Title = title
	advertisement.Name = name
	entity, err := advertisementService.Update(advertisement)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *AdvertisementController)Delete(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := advertisementService.Delete(id); err != nil {
		common.Render(ctx, "000002", nil)
	}
	common.Render(ctx, "27070400", nil)
}

func (self *AdvertisementController)BatchUpdateOrder(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
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

	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}

	advertisement.Name = name
	entity, err := advertisementService.Update(advertisement)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *AdvertisementController)SaveImage(ctx *iris.Context) {
	advertisementService := public.AdvertisementService{}
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

	advertisement, err := advertisementService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}

	advertisement.Name = name
	entity, err := advertisementService.Update(advertisement)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27070500", entity)
}
