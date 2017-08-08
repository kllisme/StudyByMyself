package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"github.com/go-errors/errors"
	"maizuo.com/soda/erp/api/src/server/service/soda"
)

type DailyBillController struct {
}

func (slef *DailyBillController) ListByBillId(ctx *iris.Context) {
	dailyBillService := &service.DailyBillService{}
	userService := &service.UserService{}
	billId := ctx.Param("id") // 账单 ID
	limit, _ := ctx.URLParamInt("limit") // 返回最多数量:Default: 10
	offset, _ := ctx.URLParamInt("offset")// 列表起始位: Default: 0
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
	list := make([]interface{},0)
	dailyBillList, err := dailyBillService.ListByBillId(limit, offset, billId)
	if err != nil {
		common.Render(ctx,"27080103",err)
		return
	}
	for _,dailyBill := range dailyBillList {
		user,err := userService.GetById(dailyBill.UserId)
		if err != nil {
			common.Logger.Debugln("获取日账单用户信息失败")
			common.Render(ctx,"27080104",err)
			return
		}
		list = append(list,dailyBill.Mapping(user))
	}
	common.Render(ctx,"27080100",&entity.PaginationData{
		Pagination:entity.Pagination{Total:total,From:offset,To:offset+limit},
		Objects:list,
	})
	return
}

func (self *DailyBillController)DetailsById(ctx *iris.Context){
	dailyBillService := service.DailyBillService{}
	userService := &service.UserService{}
	id,_ := ctx.ParamInt("id") // 日账单 ID
	limit, _ := ctx.URLParamInt("limit") // 返回最多数量:Default: 10
	offset, _ := ctx.URLParamInt("offset")// 列表起始位: Default: 0
	if id == -1 {
		common.Logger.Debugln("error id")
		common.Render(ctx,"CODE",errors.New("不合法的日账单ID"))
	}
	if limit == -1 {
		limit = 10
	}
	if offset == -1 {
		offset = 0
	}
	dailyBill,err := dailyBillService.BasicById(id)
	if err != nil {
		common.Render(ctx,"获取账单详情失败",err)
	}
	//user,err := userService.GetById(dailyBill.UserId)
	//if err != nil {
	//	common.Render(ctx,"获取账单用户失败",err)
	//}
	ticketService := soda.TicketService{}
	dailyBill.MappingDetails(nil,nil)
}
