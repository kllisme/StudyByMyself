package crm

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/kit/excel"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
	sodaModel "maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/model"
	payload "maizuo.com/soda/erp/api/src/server/payload/crm"
	"time"
	"strconv"
)

type ConsumptionController struct {
}

func (self *ConsumptionController) Paging(ctx *iris.Context) {
	userService := service.UserService{}
	ticketService := sodaService.TicketService{}
	deviceService := service.DeviceService{}
	paymentService := sodaService.PaymentService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	startAt := ctx.URLParam("startAt")         // 申请时间
	endAt := ctx.URLParam("endAt")             // 结算时间
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称
	cusMobile := ctx.URLParam("customerMobile")
	serial := ctx.URLParam("deviceSerial")

	userList := make([]*model.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if _p, err := userService.Paging(keywords, "", 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05010101", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
		if _p, err := userService.Paging("", keywords, 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05010102", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
		if len(userList) != 0 {
			for _, user := range userList {
				userIDs = append(userIDs, user.ID)
			}
		} else {
			userIDs = []int{-1}
		}
	}

	pagination, err := ticketService.Paging([]int{}, cusMobile, 0, serial, 0, userIDs, 0, 0, []int{4, 7}, startAt, endAt, offset, limit)
	if err != nil {
		common.Render(ctx, "05010103", err)
		return
	}

	ticketList := pagination.Objects.([]*sodaModel.Ticket)
	consumptionList := make([]*payload.Consumption, 0)
	paymentList, err := paymentService.GetAll()
	if err != nil {
		common.Render(ctx, "05010104", err)
		return
	}
	for _, ticket := range ticketList {
		user, err := userService.GetById(ticket.OwnerId)
		if err != nil {
			common.Render(ctx, "05010105", err)
			return
		}
		parentUser, err := userService.GetById(user.ParentID)
		if err != nil {
			common.Render(ctx, "05010106", err)
			return
		}
		device, err := deviceService.GetBySerialNumber(ticket.DeviceSerial)
		if err != nil {
			common.Render(ctx, "05010107", err)
			return
		}
		consumption := payload.Consumption{}
		consumption.TicketID = ticket.TicketId
		consumption.Agency = ""
		consumption.Mobile = user.Mobile
		consumption.ParentOperator = parentUser.Name
		consumption.ParentOperatorMobile = parentUser.Mobile
		consumption.Operator = user.Name
		consumption.Telephone = user.Telephone
		consumption.DeviceSerial = ticket.DeviceSerial
		consumption.Address = device.Address
		consumption.CustomerMobile = ticket.Mobile
		consumption.Password = ticket.Token
		consumption.Status = ticket.Status
		switch ticket.DeviceMode {
		case 1:
			consumption.TypeName = device.FirstPulseName
		case 2:
			consumption.TypeName = device.SecondPulseName
		case 3:
			consumption.TypeName = device.ThirdPulseName
		case 4:
			consumption.TypeName = device.FourthPulseName
		default:
			consumption.TypeName = "错误的数据"
		}
		consumption.Value = ticket.Value
		for _, payment := range *paymentList {
			if payment.ID == ticket.PaymentId {
				consumption.Payment = payment.Name
				break
			}
		}
		consumption.PaymentID = ticket.PaymentId
		consumption.CreatedAt = ticket.CreatedAt
		consumptionList = append(consumptionList, &consumption)
	}

	pagination.Objects = consumptionList
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
	userService := service.UserService{}
	ticketService := sodaService.TicketService{}
	deviceService := service.DeviceService{}
	paymentService := sodaService.PaymentService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	startAt := ctx.URLParam("startAt")         // 申请时间
	endAt := ctx.URLParam("endAt")             // 结算时间
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称
	cusMobile := ctx.URLParam("customerMobile")
	serial := ctx.URLParam("deviceSerial")

	userList := make([]*model.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if _p, err := userService.Paging(keywords, "", 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05010301", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
		if _p, err := userService.Paging("", keywords, 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05010302", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
		if len(userList) != 0 {
			for _, user := range userList {
				userIDs = append(userIDs, user.ID)
			}
		} else {
			userIDs = []int{-1}
		}
	}

	pagination, err := ticketService.Paging([]int{}, cusMobile, 0, serial, 0, userIDs, 0, 0, []int{4, 7}, startAt, endAt, offset, limit)
	if err != nil {
		common.Render(ctx, "05010303", err)
		return
	}

	paymentList, err := paymentService.GetAll()
	if err != nil {
		if err != nil {
			common.Render(ctx, "05010304", err)
			return
		}
	}

	ticketList := pagination.Objects.([]*sodaModel.Ticket)
	userMap := make(map[int]*model.User)
	userPage, err := userService.Paging("","", 0, 0, 0, 0)
	if err != nil {
		common.Render(ctx, "05010305", err)
		return
	}
	userList, _ = userPage.Objects.([]*model.User)
	for _, user := range userList {
		userMap[user.ID] = user
	}

	deviceMap := make(map[string]*model.Device)
	devicePage, err := deviceService.Paging(userIDs,[]int{},"",serial, 0, 0,[]int{}, 0, 0)
	if err != nil {
		common.Render(ctx, "05010306", err)
		return
	}
	deviceList, _ := devicePage.Objects.([]*model.Device)
	for _, device := range deviceList {
		deviceMap[device.SerialNumber] = device
	}

	consumptionList := make([]*payload.Consumption, 0)

	for _, ticket := range ticketList {
		user := userMap[ticket.OwnerId]
		if user == nil {
			common.Render(ctx, "05010307", err)
			return
		}
		parentUser := userMap[user.ParentID]
		if parentUser == nil {
			common.Render(ctx, "05010308", err)
			return
		}
		device := deviceMap[ticket.DeviceSerial]
		if device == nil {
			common.Render(ctx, "05010309", err)
			return
		}
		consumption := payload.Consumption{}
		consumption.TicketID = ticket.TicketId
		consumption.Agency = ""
		consumption.ParentOperator = parentUser.Name
		consumption.Mobile = user.Mobile
		consumption.ParentOperatorMobile = parentUser.Mobile
		consumption.Operator = user.Name
		consumption.Telephone = user.Telephone
		consumption.DeviceSerial = ticket.DeviceSerial
		consumption.Address = device.Address
		consumption.CustomerMobile = ticket.Mobile
		consumption.Password = ticket.Token
		switch ticket.DeviceMode {
		case 1:
			consumption.TypeName = device.FirstPulseName
		case 2:
			consumption.TypeName = device.SecondPulseName
		case 3:
			consumption.TypeName = device.ThirdPulseName
		case 4:
			consumption.TypeName = device.FourthPulseName
		default:
			consumption.TypeName = "错误的数据"
		}
		consumption.Value = ticket.Value
		for _, payment := range *paymentList {
			if payment.ID == ticket.PaymentId {
				consumption.Payment = payment.Name
				break
			}
		}
		consumption.CreatedAt = ticket.CreatedAt
		consumptionList = append(consumptionList, &consumption)
	}




	// 开始excel文件操作
	tableHead := []interface{}{"订单号", "上级运营商", "运营商名称", "服务电话", "模块编号", "楼道信息", "消费手机号", "消费密码", "类型", "消费金额", "支付方式", "下单时间"}
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
