package crm

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/kit/excel"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/common"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
	sodaModel "maizuo.com/soda/erp/api/src/server/model/soda"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	payload "maizuo.com/soda/erp/api/src/server/payload/crm"
	"time"
	"strconv"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type ConsumptionController struct {
}

func query(ctx *iris.Context) *entity.PaginationData {
	userService := mngService.UserService{}
	ticketService := sodaService.TicketService{}
	deviceService := mngService.DeviceService{}
	paymentService := sodaService.PaymentService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	startAt := ctx.URLParam("startAt")         // 申请时间
	endAt := ctx.URLParam("endAt")             // 结算时间
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称
	cusMobile := ctx.URLParam("customerMobile")
	serial := ctx.URLParam("deviceSerial")

	if keywords == "" && cusMobile == "" && serial == "" {
		common.Render(ctx, "05010108", nil)
		return nil
	}

	start, err := time.Parse("2006-01-02", startAt)
	if err != nil {
		common.Render(ctx, "05010109", err)
		return nil
	}
	end, err := time.Parse("2006-01-02", endAt)
	if err != nil {
		common.Render(ctx, "05010110", err)
		return nil
	}
	if end.After(start.AddDate(0, 0, 31)) {
		common.Render(ctx, "05010111", err)
		return nil
	}

	_userList := make([]*mngModel.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if __userList, err := userService.GetList(keywords, "", []int{}, 0); err != nil {
			common.Render(ctx, "05010101", err)
			return nil
		} else {
			_userList = append(_userList, *__userList...)
		}
		if __userList, err := userService.GetList("", keywords, []int{}, 0); err != nil {
			common.Render(ctx, "05010102", err)
			return nil
		} else {
			_userList = append(_userList, *__userList...)
		}
		if len(_userList) != 0 {
			for _, user := range _userList {
				userIDs = append(userIDs, user.ID)
			}
		} else {
			userIDs = []int{-1}
		}
		userIDs = functions.Uniq(userIDs)
	}

	pagination, err := ticketService.Paging([]int{}, cusMobile, 0, serial, 0, userIDs, 0, 0, []int{4, 7}, startAt, endAt, offset, limit)
	if err != nil {
		common.Render(ctx, "05010103", err)
		return nil
	}
	if pagination.Pagination.Total != 0 {

		ticketList := pagination.Objects.([]*sodaModel.Ticket)
		consumptionList := make([]*payload.Consumption, 0)
		paymentList, err := paymentService.GetAll()
		if err != nil {
			common.Render(ctx, "05010104", err)
			return nil
		}

		_userIDs := make([]int, 0)
		for _, ticket := range ticketList {
			_userIDs = append(_userIDs, ticket.OwnerId)
		}
		_userIDs = functions.Uniq(_userIDs)

		userList, err := userService.GetList("", "", _userIDs, 0)
		if err != nil {
			common.Render(ctx, "05010105", err)
			return nil
		}

		_parentUserIDs := make([]int, 0)
		for _, user := range *userList {
			_parentUserIDs = append(_parentUserIDs, user.ParentID)
		}
		_parentUserIDs = functions.Uniq(_parentUserIDs)

		parentUserList, err := userService.GetList("", "", _parentUserIDs, 0)
		if err != nil {
			common.Render(ctx, "05010106", err)
			return nil
		}

		_deviceSerials := make([]string, 0)
		for _, ticket := range ticketList {
			_deviceSerials = append(_deviceSerials, ticket.DeviceSerial)
		}
		_deviceSerials = functions.UniqString(_deviceSerials)

		deviceList, err := deviceService.ListBySerialNumbers(_deviceSerials...)
		if err != nil {
			common.Render(ctx, "05010107", err)
			return nil
		}

		for _, ticket := range ticketList {

			user := &mngModel.User{}
			parentUser := &mngModel.User{}
			device := &mngModel.Device{}
			consumption := payload.Consumption{}

			for _, u := range *userList {
				if u.ID == ticket.OwnerId {
					user = u
				}
			}

			if user.ID != 0 {
				for _, pU := range *parentUserList {
					if user.ParentID == pU.ID {
						parentUser = pU
					}
				}
			}

			for _, d := range *deviceList {
				if ticket.DeviceSerial == d.SerialNumber {
					device = d
				}
			}

			consumption.Map(*ticket, *user, *parentUser, *device, *paymentList)
			//common.Logger.Debug(*ticket, user, parentUser, device, *paymentList)

			consumptionList = append(consumptionList, &consumption)
		}
		pagination.Objects = consumptionList
	}
	return pagination
}

func (self *ConsumptionController) Paging(ctx *iris.Context) {

	pagination := query(ctx)
	if pagination == nil {
		return
	}
	common.Render(ctx, "05010100", pagination)
	return
}

func (self *ConsumptionController) Refund(ctx *iris.Context) {
	ticketService := sodaService.TicketService{}
	ticketId := ctx.Param("ticketId")
	if ticketId == "" {
		common.Render(ctx, "05010201", nil)
		return
	}
	ticket, err := ticketService.GetByTicketID(ticketId)
	if err != nil {
		common.Render(ctx, "05010202", err)
		return
	}
	if ticket.Status != 7 {
		common.Render(ctx, "05010203", nil)
		return
	}
	if ticket.PaymentId == 4 {
		common.Render(ctx, "05010205", nil)
		return
	}
	ticket, err = ticketService.Refund(ticketId)
	if err != nil {
		common.Render(ctx, "05010204", err)
		return
	}
	common.Render(ctx, "05010200", ticket)
}

func (self *ConsumptionController) Export(ctx *iris.Context) {
	pagination := query(ctx)
	if pagination == nil {
		return
	}
	consumptionList := pagination.Objects.([]*payload.Consumption)
	// 开始excel文件操作
	tableHead := []interface{}{"订单号", "上级运营商", "运营商名称", "服务电话", "模块编号", "楼道信息", "消费手机号", "消费密码", "类型", "消费金额", "支付方式", "下单时间", "状态"}
	tableName := "消费查询列表"

	fileName := "消费订单详情" + time.Now().Format("20060102") + strconv.FormatInt(time.Now().Local().Unix(), 10)

	sheet, file, fileUrl, fileName, err := excel.GetExcelHeader(fileName, tableHead, tableName)
	if err != nil {
		common.Logger.Warningln("操作excel文件失败, err ------------>", err)
		common.Render(ctx, "05010310", err)
		return
	}
	//将查询的数据装填
	for _, consumption := range consumptionList {
		if excel.ExportConsumptionAsCol(sheet, consumption) == 0 {
			common.Logger.Warningln("excel文件插入记录失败,err ------------>", err)
			common.Render(ctx, "05010311", err)
			return
		}
	}
	err = file.Save(fileUrl)
	if err != nil {
		common.Logger.Warningln("excel文件保存失败,err ------------>", err)
		common.Render(ctx, "05010312", err)
		return
	}
	sendFile := viper.GetString("server.href") + viper.GetString("export.loadsPath") + "/" + fileName
	common.Render(ctx, "05010300", map[string]string{"url": sendFile})
	return
}
