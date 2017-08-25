package public

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	"maizuo.com/soda/erp/api/src/server/common"
)

type RegionController struct {

}

func (self *RegionController)GetByID(ctx *iris.Context) {
	regionService := public.RegionService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04010102", err)
		return
	}
	region, err := regionService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04010101", err)
	}
	common.Render(ctx, "04010100", region)
}

func (self *RegionController)GetCities(ctx *iris.Context) {
	regionService := public.RegionService{}
	provinceID, err := ctx.ParamInt("province_id")
	if err != nil {
		common.Render(ctx, "04010301", err)
		return
	}
	province, err := regionService.GetByID(provinceID)
	if err != nil {
		common.Render(ctx, "04010303", err)
		return
	}
	region, err := regionService.GetCities(province.ID)
	if err != nil {
		common.Render(ctx, "04010302", err)
	}
	common.Render(ctx, "04010300", region)
}

func (self *RegionController)GetRegions(ctx *iris.Context) {
	regionService := public.RegionService{}
	cityID, err := ctx.ParamInt("city_id")
	if err != nil {
		common.Render(ctx, "04010401", err)
		return
	}
	city,err:=regionService.GetByID(cityID)
	if err != nil {
		common.Render(ctx, "04010403", err)
		return
	}
	region, err := regionService.GetRegions(city.ID)
	if err != nil {
		common.Render(ctx, "04010402", err)
	}
	common.Render(ctx, "04010400", region)
}


func (self *RegionController)GetProvinces(ctx *iris.Context) {
	regionService := public.RegionService{}
	region, err := regionService.GetProvinces()
	if err != nil {
		common.Render(ctx, "04010201", err)
	}
	common.Render(ctx, "04010200", region)
}
//
//func (self *RegionController)Paging(ctx *iris.Context) {
//	regionService := public.RegionService{}
//	page, _ := ctx.URLParamInt("page")
//	perPage, _ := ctx.URLParamInt("per_page")
//	result, err := regionService.Paging(page, perPage)
//	if err != nil {
//		common.Render(ctx, "000002", nil)
//		return
//	}
//	common.Render(ctx, "27040200", result)
//	return
//}
