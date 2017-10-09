package crm

import (
	"time"
	sodaModel "maizuo.com/soda/erp/api/src/server/model/soda"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/service/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	sodaService "maizuo.com/soda/erp/api/src/server/service/soda"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	"errors"
	"github.com/jinzhu/gorm"
)

type Customer struct {
	ID               int        `json:"id"`
	NickName         string        `json:"nickName"`
	Mobile           string        `json:"mobile"`
	School           string        `json:"school"`
	WechatName       string        `json:"wechatName"`
	CreatedAt        time.Time        `json:"createdAt"`
	OpenID           string        `json:"openId"`
	WalletCount      int        `json:"walletCount"`
	ChipcardCount    int        `json:"chipcardCount"`
	RecentAddress    string        `json:"recentAddress"`
	LastTicketResume string        `json:"lastTicketResume"`
	LastTicketID     string        `json:"lastTicketId"`
}

func (self *Customer) Map(user *sodaModel.User) error {
	applicationService := soda_manager.ApplicationService{}
	//provinceService := soda_manager.ProvinceService{}
	//cityService := soda_manager.CityService{}
	//areaService := soda_manager.AreaService{}
	schoolService := soda_manager.SchoolService{}
	walletService := sodaService.WalletService{}
	chipcardService := sodaService.ChipcardService{}
	ticketService := sodaService.TicketService{}
	deviceService := mngService.DeviceService{}
	wallet, err := walletService.GetByUserMobile(user.Mobile)
	if err != nil {
		common.Logger.Error("钱包信息出错")
		return errors.New("钱包信息出错")
	}
	chipcard, err := chipcardService.GetByUserMobile(user.Mobile)
	if err != nil {
		if err == gorm.ErrRecordNotFound{
			self.ChipcardCount = -1
		}else {
			common.Logger.Error("芯片卡信息出错")
			return errors.New("芯片卡信息出错")
		}
	} else {
		self.ChipcardCount = chipcard.Value
	}
	ticket := sodaModel.Ticket{}
	ticketPage, err := ticketService.Paging([]int{}, user.Mobile, 0, "", 0, []int{}, 0, 0, []int{}, "", "", 0, 1)
	if err != nil {
		common.Logger.Error("票据信息出错")
	}
	if ticketPage.Pagination.Total != 0 {
		ticketList := ticketPage.Objects.([]*sodaModel.Ticket)
		ticket = *ticketList[0]
	}
	if ticket.ID != 0 {
		d, err := deviceService.GetBySerialNumber(ticket.DeviceSerial)
		if err != nil {
			common.Logger.Error("设备信息出错")
		}
		device := *d

		self.LastTicketID = ticket.TicketId
		application, err := applicationService.GetByID(ticket.APPID)
		if err != nil {
			common.Logger.Error("应用信息出错")
			//return errors.New("应用信息出错")
			application = &mngModel.Application{}
		}
		typeName := ""
		switch ticket.DeviceMode {
		case 1:
			typeName = device.FirstPulseName
		case 2:
			typeName = device.SecondPulseName
		case 3:
			typeName = device.ThirdPulseName
		case 4:
			typeName = device.FourthPulseName
		default:
			typeName = "/"
		}
		value := functions.Float64ToString(float64(ticket.Value) / 100.00, 2)
		self.LastTicketResume = application.Name + "-" + typeName + device.SerialNumber + "-" + value + "元" + " (密码：" + ticket.Token + ")"

		//if device.ProvinceID != 0 {
		//	p, err := provinceService.GetByCode(device.ProvinceID)
		//	if err != nil {
		//		common.Logger.Error("省信息出错")
		//	} else {
		//		self.RecentAddress += p.Name
		//	}
		//}
		//
		//if device.CityID != 0 {
		//	c, err := cityService.GetByCode(device.CityID)
		//	if err != nil {
		//		common.Logger.Error("市信息出错")
		//	} else {
		//		self.RecentAddress += c.Name
		//	}
		//}
		//
		//if device.DistrictID != 0 {
		//	a, err := areaService.GetByCode(device.DistrictID)
		//	if err != nil {
		//		common.Logger.Error("区信息出错")
		//	} else {
		//		self.RecentAddress += a.Name
		//	}
		//}

		if device.SchoolID != 0 {
			s, err := schoolService.GetByID(device.SchoolID)
			if err != nil {
				common.Logger.Error("学校信息出错")
			} else {
				self.RecentAddress += s.Name
			}
		}

		self.RecentAddress += device.Address
	}

	self.ID = user.ID
	self.Mobile = user.Mobile
	self.NickName = user.NickName
	self.WechatName = "无"
	self.CreatedAt = user.CreatedAt
	self.OpenID = "无"
	self.WalletCount = wallet.Value
	return nil
}
