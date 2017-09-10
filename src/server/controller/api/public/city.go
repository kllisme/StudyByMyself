package public

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/public"
	model "maizuo.com/soda/erp/api/src/server/model/public"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
)

type CityController struct {

}

func (self *CityController)GetByID(ctx *iris.Context) {
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04050101", err)
		return
	}
	city, err := cityService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04050102", err)
		return
	}
	if city.ParentCode != "" {
		province, err := provinceService.GetByCode(city.ParentCode)
		if err != nil {
			common.Render(ctx, "04050103", err)
			return
		}
		city.ProvinceCode = city.ParentCode
		city.ProvinceName = province.Name
	}

	common.Render(ctx, "04050100", city)
}

func (self *CityController)Paging(ctx *iris.Context) {
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	provinceCode := strings.TrimSpace(ctx.URLParam("provinceCode"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	result, err := cityService.Paging(name, provinceCode, offset, limit)
	if err != nil {
		common.Render(ctx, "04050201", err)
		return
	}

	cityList := result.Objects.([]*model.City)
	for _, city := range cityList {
		if city.ParentCode != "" {
			province, err := provinceService.GetByCode(city.ParentCode)
			if err != nil {
				common.Render(ctx, "04050202", err)
				return
			}
			city.ProvinceCode = city.ParentCode
			city.ProvinceName = province.Name
		}
	}

	common.Render(ctx, "04050200", result)
	return
}

func (self *CityController)Create(ctx *iris.Context) {
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04050301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04050302", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "04050303", nil)
		return
	}
	if _, err := cityService.GetByCode(code); err != gorm.ErrRecordNotFound {
		common.Render(ctx, "04050304", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode == "" {
		common.Render(ctx, "04050305", nil)
		return
	} else if _, err := provinceService.GetByCode(provinceCode); err != nil {
		common.Render(ctx, "04050306", err)
		return
	}
	city := model.City{
		Name:name,
		Code:code,
		ParentCode:provinceCode,
	}
	entity, err := cityService.Create(&city)
	if err != nil {
		common.Render(ctx, "04050307", err)
		return
	}
	common.Render(ctx, "04050300", entity)
}

func (self *CityController)Update(ctx *iris.Context) {
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

	city, err := cityService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27070502", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "27070303", nil)
		return
	}

	if _, err := cityService.GetByCode(code); err != nil {
		if err != gorm.ErrRecordNotFound {
			common.Render(ctx, "27070303", nil)
			return
		}

	} else if city.Code != code {
		common.Render(ctx, "27070303", nil)
		return
	}
	provinceCode := strings.TrimSpace(params.Get("provinceCode").MustString())
	if provinceCode == "" {
		common.Render(ctx, "27070305", nil)
		return
	} else if _, err := provinceService.GetByCode(provinceCode); err != nil {
		common.Render(ctx, "27070306", err)
		return
	}
	city.Name = name
	city.Code = code
	city.ParentCode = provinceCode
	entity, err := cityService.Update(city)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *CityController)Delete(ctx *iris.Context) {
	cityService := public.CityService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := cityService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070400", nil)
}
