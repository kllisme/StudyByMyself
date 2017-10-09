package crm

import (
	"github.com/bitly/go-simplejson"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
	"maizuo.com/soda/erp/api/src/server/payload/crm"
	"strings"
)

type CustomerController struct{}

func (self *CustomerController) GetByMobile(ctx *iris.Context) {
	userService := sodaService.UserService{}

	mobile := ctx.Param("mobile")
	if len(mobile) != 11 {
		common.Render(ctx, "05030101", nil)
		return
	}
	customer := crm.Customer{}
	userEntity, err := userService.GetByMobile(mobile)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			common.Render(ctx, "05040102", err)
			return
		}
		common.Render(ctx, "05040103", err)
		return
	}
	err = customer.Map(userEntity)
	if err != nil {
		common.Render(ctx, "05040104", err)
		return
	}
	common.Render(ctx, "05040100", customer)
}

func (self *CustomerController) ChangePassword(ctx *iris.Context) {
	userService := sodaService.UserService{}
	mobile := ctx.Param("mobile")

	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "05040201", err)
		return
	}
	newPassword := strings.TrimSpace(params.Get("newPassword").MustString())
	if len(newPassword) < 6 || len(newPassword) > 16 {
		common.Render(ctx, "05040202", nil)
		return
	}
	entity, err := userService.ChangePassword(mobile, newPassword)
	if err != nil {
		common.Render(ctx, "05040203", err)
		return
	}
	common.Render(ctx, "05040200", entity)
}
