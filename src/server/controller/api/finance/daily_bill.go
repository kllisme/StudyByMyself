package finance

import (
	"github.com/go-errors/errors"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
)

type DailyBillController struct {
}

func (slef *DailyBillController) ListByBillId(ctx *iris.Context) {
	dailyBillService := &mngService.DailyBillService{}
	userService := &mngService.UserService{}
	id := ctx.Param("id")                  // 账单 ID 不是真正的id,而是bill中的bill_id
	limit, _ := ctx.URLParamInt("limit")   // 返回最多数量:Default: 10
	offset, _ := ctx.URLParamInt("offset") // 列表起始位: Default: 0
	if id == "" {
		common.Render(ctx, "27090105", errors.New("billId为空"))
	}
	total, err := dailyBillService.TotalByBillId(id)
	if err != nil {
		common.Render(ctx, "27090101", err)
		return
	}
	if total == 0 {
		common.Render(ctx, "27090102", nil)
		return
	}
	if limit == 0 {
		limit = 10
	}
	objects := make([]interface{}, 0)
	dailyBillList, err := dailyBillService.ListByBillId(limit, offset, id)
	if err != nil {
		common.Render(ctx, "27090103", err)
		return
	}
	for _, dailyBill := range dailyBillList {
		user, err := userService.GetById(dailyBill.UserId)
		if err != nil {
			common.Logger.Debugln("获取日账单用户信息失败")
			common.Render(ctx, "27090104", err)
			return
		}
		objects = append(objects, dailyBill.Mapping(user))
	}
	common.Render(ctx, "27090100", &entity.PaginationData{
		Pagination: entity.Pagination{Total: total, From: offset, To: offset + limit - 1},
		Objects:    objects,
	})
	return
}

func (self *DailyBillController) DetailsById(ctx *iris.Context) {
	dailyBillService := &mngService.DailyBillService{}
	ticketService := &sodaService.TicketService{}
	deviceService := &mngService.DeviceService{}
	id, _ := ctx.ParamInt("id")            // 日账单 ID
	limit, _ := ctx.URLParamInt("limit")   // 返回最多数量:Default: 10
	offset, _ := ctx.URLParamInt("offset") // 列表起始位: Default: 0
	if id == 0 {
		common.Logger.Debugln("error id", id)
		common.Render(ctx, "27090201", nil)
		return
	}
	if limit == 0 {
		limit = 10
	}
	dailyBill, err := dailyBillService.BasicById(id)
	if err != nil {
		common.Render(ctx, "27090202", err)
		return
	}
	total, err := ticketService.TotalByDailyBill(dailyBill)
	if err != nil {
		common.Render(ctx, "27090203", err)
		return
	}
	tickets, err := ticketService.DetailsByDailyBill(dailyBill, limit, offset)
	if err != nil {
		common.Render(ctx, "27090204", err)
		return
	}
	objects := make([]interface{}, 0)
	for _, ticket := range tickets {
		device, err := deviceService.GetBySerialNumber(ticket.DeviceSerial)
		if err != nil {
			common.Render(ctx, "27090205", err)
			return
		}
		objects = append(objects, ticket.Mapping(device, dailyBill))
	}
	common.Render(ctx, "27090200", &entity.PaginationData{
		Pagination: entity.Pagination{Total: total, From: offset, To: offset + limit - 1},
		Objects:    objects,
	})
	return
}
