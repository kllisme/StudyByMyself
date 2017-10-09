package crm

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
)

type BillController struct {

}

func (self *BillController) Paging(ctx *iris.Context) {
	billService := sodaService.BillService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	startAt := ctx.URLParam("startAt")         // 申请时间
	endAt := ctx.URLParam("endAt")             // 结算时间
	_type, _ := ctx.URLParamInt("type")               // 运营商名称、帐号名称
	cusMobile := ctx.URLParam("customerMobile")
	action, _ := ctx.URLParamInt("action")

	if cusMobile == "" {
		common.Render(ctx, "05030101", nil)
		return
	}

	//start, err := time.Parse("2006-01-02", startAt)
	//if err != nil {
	//	common.Render(ctx, "05010109", err)
	//	return nil
	//}
	//end, err := time.Parse("2006-01-02", endAt)
	//if err != nil {
	//	common.Render(ctx, "05010110", err)
	//	return nil
	//}
	//if end.After(start.AddDate(0, 0, 31)) {
	//	common.Render(ctx, "05010111", err)
	//	return nil
	//}

	pagination, err := billService.Paging(cusMobile, action, _type, startAt, endAt, offset, limit)
	if err != nil {
		common.Render(ctx, "05030102", err)
		return
	}
	common.Render(ctx, "05030100", pagination)
	return
}

func (self *BillController) PagingChipcardBill(ctx *iris.Context) {
	chipcardBillService := sodaService.ChipcardBillService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	startAt := ctx.URLParam("startAt")         // 申请时间
	endAt := ctx.URLParam("endAt")             // 结算时间
	cusMobile := ctx.URLParam("customerMobile")
	action, _ := ctx.URLParamInt("action")

	if cusMobile == "" {
		common.Render(ctx, "05030201", nil)
		return
	}

	//start, err := time.Parse("2006-01-02", startAt)
	//if err != nil {
	//	common.Render(ctx, "05010109", err)
	//	return nil
	//}
	//end, err := time.Parse("2006-01-02", endAt)
	//if err != nil {
	//	common.Render(ctx, "05010110", err)
	//	return nil
	//}
	//if end.After(start.AddDate(0, 0, 31)) {
	//	common.Render(ctx, "05010111", err)
	//	return nil
	//}

	pagination, err := chipcardBillService.Paging(cusMobile, action, 0, startAt, endAt, offset, limit)
	if err != nil {
		common.Render(ctx, "05030202", err)
		return
	}
	common.Render(ctx, "05030200", pagination)
	return
}
