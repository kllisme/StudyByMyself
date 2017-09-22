package soda_manager

import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	model "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
)

type AreaController struct {

}

func (self *AreaController)GetByID(ctx *iris.Context) {
	areaService := mngService.AreaService{}
	cityService := mngService.CityService{}
	provinceService:= mngService.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04070101", err)
		return
	}
	area, err := areaService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04070102", err)
		return
	}
	if area.ParentCode != "" {
		city, err := cityService.GetByCode(area.ParentCode)
		if err != nil {
			common.Render(ctx, "04070103", err)
			return
		}
		area.CityCode = area.ParentCode
		area.CityName = city.Name
		if city.ParentCode != "" {
			province, err := provinceService.GetByCode(city.ParentCode)
			if err != nil {
				common.Render(ctx, "04070104", err)
				return
			}
			area.ProvinceCode = city.ParentCode
			area.ProvinceName = province.Name
		}
	}

	common.Render(ctx, "04070100", area)
}

func (self *AreaController)Paging(ctx *iris.Context) {
	areaService := mngService.AreaService{}
	cityService := mngService.CityService{}
	provinceService:= mngService.ProvinceService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	cityCode := strings.TrimSpace(ctx.URLParam("cityCode"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	result, err := areaService.Paging(name, cityCode, offset, limit)
	if err != nil {
		common.Render(ctx, "04070201", err)
		return
	}

	areaList := result.Objects.([]*model.Area)
	for _, area := range areaList {
		if area.ParentCode != "" {
			city, err := cityService.GetByCode(area.ParentCode)
			if err != nil {
				common.Render(ctx, "04070202", err)
				return
			}
			area.CityCode = area.ParentCode
			area.CityName = city.Name
			if city.ParentCode != "" {
				province, err := provinceService.GetByCode(city.ParentCode)
				if err != nil {
					common.Render(ctx, "04070203", err)
					return
				}
				area.ProvinceCode = city.ParentCode
				area.ProvinceName = province.Name
			}
		}
	}

	common.Render(ctx, "04070200", result)
	return
}

func (self *AreaController)Create(ctx *iris.Context) {
	areaService := mngService.AreaService{}
	cityService := mngService.CityService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04070301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04070302", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "04070303", nil)
		return
	}
	if _, err := areaService.GetByCode(code); err != gorm.ErrRecordNotFound {
		common.Render(ctx, "04070304", nil)
		return
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode == "" {
		common.Render(ctx, "04070305", nil)
		return
	} else if _, err := cityService.GetByCode(cityCode); err != nil {
		common.Render(ctx, "04070306", err)
		return
	}
	area := model.Area{
		Name:name,
		Code:code,
		ParentCode:cityCode,
	}
	entity, err := areaService.Create(&area)
	if err != nil {
		common.Render(ctx, "04070307", err)
		return
	}
	common.Render(ctx, "04070300", entity)
}

func (self *AreaController)Update(ctx *iris.Context) {
	areaService := mngService.AreaService{}
	cityService := mngService.CityService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04070401", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04070402", err)
		return
	}

	area, err := areaService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04070403", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04070404", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "04070405", nil)
		return
	}

	if _, err := areaService.GetByCode(code); err != nil {
		if err != gorm.ErrRecordNotFound {
			common.Render(ctx, "04070406", nil)
			return
		}

	} else if area.Code != code {
		common.Render(ctx, "04070407", nil)
		return
	}
	cityCode := strings.TrimSpace(params.Get("cityCode").MustString())
	if cityCode == "" {
		common.Render(ctx, "04070408", nil)
		return
	} else if _, err := cityService.GetByCode(cityCode); err != nil {
		common.Render(ctx, "04070409", err)
		return
	}
	area.Name = name
	area.Code = code
	area.ParentCode = cityCode
	entity, err := areaService.Update(area)
	if err != nil {
		common.Render(ctx, "04070410", err)
		return
	}
	common.Render(ctx, "04070401", entity)
}

func (self *AreaController)Delete(ctx *iris.Context) {
	areaService := mngService.AreaService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04070501", err)
		return
	}
	if err := areaService.Delete(id); err != nil {
		common.Render(ctx, "04070502", err)
		return
	}
	common.Render(ctx, "04070500", nil)
}
