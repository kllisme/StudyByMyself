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

type ProvinceController struct {

}

func (self *ProvinceController)GetByID(ctx *iris.Context) {
	provinceService := public.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04060101", err)
		return
	}
	province, err := provinceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04060102", err)
		return
	}
	common.Render(ctx, "04060100", province)
}

func (self *ProvinceController)Paging(ctx *iris.Context) {
	provinceService := public.ProvinceService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := provinceService.Paging(offset, limit)
	if err != nil {
		common.Render(ctx, "04060201", err)
		return
	}
	common.Render(ctx, "04060200", result)
	return
}

func (self *ProvinceController)Create(ctx *iris.Context) {
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04060301", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04060302", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "04060303", nil)
		return
	}
	if _, err := provinceService.GetByCode(code); err != gorm.ErrRecordNotFound {
		common.Render(ctx, "04060304", nil)
		return
	}
	province := model.Province{
		Name:name,
		Code:code,
	}
	entity, err := provinceService.Create(&province)
	if err != nil {
		common.Render(ctx, "04060305", err)
		return
	}
	common.Render(ctx, "04060300", entity)
}

func (self *ProvinceController)Update(ctx *iris.Context) {
	provinceService := public.ProvinceService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "04060401", err)
		return
	}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04060402", err)
		return
	}

	province, err := provinceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04060403", err)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "04060404", nil)
		return
	}
	code := strings.TrimSpace(params.Get("code").MustString())
	if code == "" {
		common.Render(ctx, "04060405", nil)
		return
	}

	if _, err := provinceService.GetByCode(code); err != nil {
		if err != gorm.ErrRecordNotFound {
			common.Render(ctx, "04060406", nil)
			return
		}

	} else if province.Code != code {
		common.Render(ctx, "04060407", nil)
		return
	}
	province.Name = name
	province.Code = code
	entity, err := provinceService.Update(province)
	if err != nil {
		common.Render(ctx, "04060408", err)
		return
	}
	common.Render(ctx, "04060400", entity)
}

func (self *ProvinceController)Delete(ctx *iris.Context) {
	provinceService := public.ProvinceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04060501", err)
		return
	}
	if err := provinceService.Delete(id); err != nil {
		common.Render(ctx, "04060502", err)
	}
	common.Render(ctx, "04060500", nil)
}
