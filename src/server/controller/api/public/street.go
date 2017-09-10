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

type StreetController struct {

}

func (self *StreetController)GetByID(ctx *iris.Context) {
	streetService := public.StreetService{}
	areaService := public.AreaService{}
	cityService := public.CityService{}
	provinceService := public.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	street, err := streetService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	if street.ParentCode != "" {
		area, err := areaService.GetByCode(street.ParentCode)
		if err != nil {
			common.Render(ctx, "000002", err)
			return
		}
		street.AreaCode = street.ParentCode
		street.AreaName = area.Name
		if area.ParentCode != "" {
			city, err := cityService.GetByCode(area.ParentCode)
			if err != nil {
				common.Render(ctx, "000002", err)
				return
			}
			street.CityCode = area.ParentCode
			street.CityName = city.Name
			if city.ParentCode != "" {
				province, err := provinceService.GetByCode(city.ParentCode)
				if err != nil {
					common.Render(ctx, "000002", err)
					return
				}
				street.ProvinceCode = city.ParentCode
				street.ProvinceName = province.Name
			}
		}
	}

	common.Render(ctx, "27070100", street)
}

func (self *StreetController)Paging(ctx *iris.Context) {
	streetService := public.StreetService{}
	areaService := public.AreaService{}
	provinceService := public.ProvinceService{}
	cityService := public.CityService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	areaCode := strings.TrimSpace(ctx.URLParam("areaCode"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	result, err := streetService.Paging(name, areaCode, offset, limit)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}

	streetList := result.Objects.([]*model.Street)
	for _, street := range streetList {
		if street.ParentCode != "" {
			area, err := areaService.GetByCode(street.ParentCode)
			if err != nil {
				common.Render(ctx, "000002", err)
				return
			}
			street.AreaCode = street.ParentCode
			street.AreaName = area.Name
			if area.ParentCode != "" {
				city, err := cityService.GetByCode(area.ParentCode)
				if err != nil {
					common.Render(ctx, "000002", err)
					return
				}
				street.CityCode = area.ParentCode
				street.CityName = city.Name
				if city.ParentCode != "" {
					province, err := provinceService.GetByCode(city.ParentCode)
					if err != nil {
						common.Render(ctx, "000002", err)
						return
					}
					street.ProvinceCode = city.ParentCode
					street.ProvinceName = province.Name
				}
			}
		}
	}

	common.Render(ctx, "27070200", result)
	return
}

func (self *StreetController)Create(ctx *iris.Context) {
	streetService := public.StreetService{}
	areaService := public.AreaService{}
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
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "27070303", nil)
		return
	}
	if _, err := streetService.GetByCode(code); err != gorm.ErrRecordNotFound {
		common.Render(ctx, "27070304", nil)
		return
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode == "" {
		common.Render(ctx, "27070305", nil)
		return
	} else if _, err := areaService.GetByCode(areaCode); err != nil {
		common.Render(ctx, "27070306", err)
		return
	}
	street := model.Street{
		Name:name,
		Code:code,
		ParentCode:areaCode,
	}
	entity, err := streetService.Create(&street)
	if err != nil {
		common.Render(ctx, "27070307", err)
		return
	}
	common.Render(ctx, "27070300", entity)
}

func (self *StreetController)Update(ctx *iris.Context) {
	streetService := public.StreetService{}
	areaService := public.AreaService{}
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

	street, err := streetService.GetByID(id)
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

	if _, err := streetService.GetByCode(code); err != nil {
		if err != gorm.ErrRecordNotFound {
			common.Render(ctx, "27070303", nil)
			return
		}

	} else if street.Code != code {
		common.Render(ctx, "27070303", nil)
		return
	}
	areaCode := strings.TrimSpace(params.Get("areaCode").MustString())
	if areaCode == "" {
		common.Render(ctx, "27070305", nil)
		return
	} else if _, err := areaService.GetByCode(areaCode); err != nil {
		common.Render(ctx, "27070306", err)
		return
	}
	street.Name = name
	street.Code = code
	street.ParentCode = areaCode
	entity, err := streetService.Update(street)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27070500", entity)
}

func (self *StreetController)Delete(ctx *iris.Context) {
	streetService := public.StreetService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := streetService.Delete(id); err != nil {
		common.Render(ctx, "000002", err)
	}
	common.Render(ctx, "27070400", nil)
}
