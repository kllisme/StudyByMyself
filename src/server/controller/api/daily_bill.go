package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
)

type DailyBillController struct {
}

func (slef *DailyBillController) ListByBillId(ctx *iris.Context) {
	dailyBillService := &service.DailyBillService{}

	limit, _ := ctx.URLParamInt("limit") // 返回最多数量:Default: 10
	offset, _ := ctx.URLParamInt("offset")// 列表起始位: Default: 0
	billId := ctx.URLParam("id") // 账单 ID
	//cashAccountType := functions.StringToInt(params["cashAccountType"]) //提现方式
	total, err := dailyBillService.TotalByBillId(billId)
	if err != nil {
		common.Render(ctx,"27080101",err)
		return
	}
	if total == 0 {
		common.Render(ctx,"27080102",err)
		return
	}
	if limit == -1 {
		limit = 10
	}
	if offset == -1 {
		offset = 0
	}
	list, err := dailyBillService.ListByBillId(limit, offset, billId)
	if err != nil {
		common.Render(ctx,"27080103",err)
		return
	}
	common.Render(ctx,"27080100",&entity.PaginationData{
		Pagination:entity.Pagination{Total:total,From:offset,To:offset+limit},
		Objects:list,
	})
	return
}
