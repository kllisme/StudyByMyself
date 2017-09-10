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
		common.Render(ctx, "000003", err)
		return
	}
	school, err := schoolService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	if school.AreaCode != "" {
		area, err := areaService.GetByCode(school.AreaCode)
		if err != nil {
			common.Render(ctx, "000002", err)
			return
		}
		school.AreaName = area.Name

	}
	if school.CityCode != "" {
		city, err := cityService.GetByCode(school.CityCode)
		if err != nil {
			common.Render(ctx, "000002", err)
			return
		}
		school.CityName = city.Name

	}
	if school.ProvinceCode != "" {
		province, err := provinceService.GetByCode(school.ProvinceCode)
		if err != nil {
			common.Render(ctx, "000002", err)
			return
		}
		school.ProvinceName = province.Name
	}

	common.Render(ctx, "27070100", school)
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
		common.Render(ctx, "000002", err)
		return
	}

	schoolList := result.Objects.([]*model.School)
	for _, school := range schoolList {
		if school.AreaCode != "" {
			area, err := areaService.GetByCode(school.AreaCode)
			if err != nil {
				common.Render(ctx, "000002", err)
				return
			}
			school.AreaName = area.Name

		}
		if school.CityCode != "" {
			city, err := cityService.GetByCode(school.CityCode)
			if err != nil {
				common.Render(ctx, "000002", err)
				return
			}
			school.CityName = city.Name

		}
		if school.ProvinceCode != "" {
			province, err := provinceService.GetByCode(school.ProvinceCode)
			if err != nil {
				common.Render(ctx, "000002", err)
				return
			}
			school.ProvinceName = province.Name
		}
	}

	common.Render(ctx, "27070200", result)
	return
}

func (self *SchoolController)Create(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070302", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode != "" {
		if _, err := provinceService.GetByCode(provinceCode); err != nil {
			common.Render(ctx, "27070304", nil)
			return
		}
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode != "" {
		if _, err := cityService.GetByCode(cityCode); err != nil {
			common.Render(ctx, "27070304", nil)
			return
		}
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode != "" {
		if _, err := areaService.GetByCode(areaCode); err != nil {
			common.Render(ctx, "27070304", nil)
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
		common.Render(ctx, "27070307", err)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *SchoolController)Update(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27070501", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	school, err := schoolService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode != "" {
		if _, err := provinceService.GetByCode(provinceCode); err != nil {
			common.Render(ctx, "27070304", nil)
			return
		}
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode != "" {
		if _, err := cityService.GetByCode(cityCode); err != nil {
			common.Render(ctx, "27070304", nil)
			return
		}
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode != "" {
		if _, err := areaService.GetByCode(areaCode); err != nil {
			common.Render(ctx, "27070304", nil)
			return
		}
	}
	school.Name = name
	school.ProvinceCode = provinceCode
	school.CityCode = cityCode
	school.AreaCode = areaCode
	entity, err := schoolService.Update(school)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *SchoolController)Delete(ctx *iris.Context) {
	schoolService := public.SchoolService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := schoolService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070400", nil)
}
