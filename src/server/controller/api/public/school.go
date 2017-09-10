package public

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	model "maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
)

type SchoolController struct {

}

func (self *SchoolController)GetByID(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04010101", err)
		return
	}
	school, err := schoolService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04010102", err)
		return
	}
	if school.AreaCode != "" {
		area, err := areaService.GetByCode(school.AreaCode)
		if err != nil {
			common.Render(ctx, "04010103", err)
			return
		}
		school.AreaName = area.Name

	}
	if school.CityCode != "" {
		city, err := cityService.GetByCode(school.CityCode)
		if err != nil {
			common.Render(ctx, "04010104", err)
			return
		}
		school.CityName = city.Name

	}
	if school.ProvinceCode != "" {
		province, err := provinceService.GetByCode(school.ProvinceCode)
		if err != nil {
			common.Render(ctx, "04010105", err)
			return
		}
		school.ProvinceName = province.Name
	}

	common.Render(ctx, "04010100", school)
}

func (self *SchoolController)Paging(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	provinceService := public.ProvinceService{}
	cityService := public.CityService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	name := strings.TrimSpace(ctx.URLParam("name"))
	areaCode := strings.TrimSpace(ctx.URLParam("areaCode"))
	cityCode := strings.TrimSpace(ctx.URLParam("cityCode"))
	provinceCode := strings.TrimSpace(ctx.URLParam("provinceCode"))

	result, err := schoolService.Paging(name, provinceCode, cityCode, areaCode, offset, limit)
	if err != nil {
		common.Render(ctx, "04010201", err)
		return
	}

	schoolList := result.Objects.([]*model.School)
	for _, school := range schoolList {
		if school.AreaCode != "" {
			area, err := areaService.GetByCode(school.AreaCode)
			if err != nil {
				common.Render(ctx, "04010202", err)
				return
			}
			school.AreaName = area.Name

		}
		if school.CityCode != "" {
			city, err := cityService.GetByCode(school.CityCode)
			if err != nil {
				common.Render(ctx, "04010203", err)
				return
			}
			school.CityName = city.Name

		}
		if school.ProvinceCode != "" {
			province, err := provinceService.GetByCode(school.ProvinceCode)
			if err != nil {
				common.Render(ctx, "04010204", err)
				return
			}
			school.ProvinceName = province.Name
		}
	}

	common.Render(ctx, "04010200", result)
	return
}

func (self *SchoolController)Create(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04010301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04010302", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode != "" {
		if _, err := provinceService.GetByCode(provinceCode); err != nil {
			common.Render(ctx, "04010303", nil)
			return
		}
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode != "" {
		if _, err := cityService.GetByCode(cityCode); err != nil {
			common.Render(ctx, "04010304", nil)
			return
		}
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode != "" {
		if _, err := areaService.GetByCode(areaCode); err != nil {
			common.Render(ctx, "04010305", nil)
			return
		}
	}
	school := model.School{
		Name:name,
		ProvinceCode:provinceCode,
		CityCode:cityCode,
		AreaCode:areaCode,
	}
	entity, err := schoolService.Create(&school)
	if err != nil {
		common.Render(ctx, "04010306", err)
		return
	}
	common.Render(ctx, "04010300", entity)
}

func (self *SchoolController)Update(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04010401", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04010402", err)
		return
	}

	school, err := schoolService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04010403", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04010404", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode != "" {
		if _, err := provinceService.GetByCode(provinceCode); err != nil {
			common.Render(ctx, "04010405", nil)
			return
		}
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode != "" {
		if _, err := cityService.GetByCode(cityCode); err != nil {
			common.Render(ctx, "04010406", nil)
			return
		}
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode != "" {
		if _, err := areaService.GetByCode(areaCode); err != nil {
			common.Render(ctx, "04010407", nil)
			return
		}
	}
	school.Name = name
	school.ProvinceCode = provinceCode
	school.CityCode = cityCode
	school.AreaCode = areaCode
	entity, err := schoolService.Update(school)
	if err != nil {
		common.Render(ctx, "04010408", err)
		return
	}
	common.Render(ctx, "04010400", entity)
}

func (self *SchoolController)Delete(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04010501", err)
		return
	}
	if err := schoolService.Delete(id); err != nil {
		common.Render(ctx, "04010502", err)
	}
	common.Render(ctx, "04010500", nil)
}
