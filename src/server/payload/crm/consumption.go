package crm

import (
	"time"
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type Consumption struct {
	TicketID             string        `json:"ticketId"`
	Agency               string        `json:"agency"`
	Operator             string        `json:"operator"`
	ParentOperator       string        `json:"parentOperator"`
	ParentOperatorMobile string        `json:"parentOperatorMobile"`
	Mobile               string        `json:"mobile"`
	Telephone            string        `json:"telephone"`
	DeviceSerial         string        `json:"deviceSerial"`
	Address              string        `json:"address"`
	CustomerMobile       string        `json:"customerMobile"`
	Password             string        `json:"password"`
	TypeName             string        `json:"typeName"`
	Value                int        `json:"value"`
	Payment              string        `json:"payment"`
	PaymentID            int        `json:"paymentId"`
	CreatedAt            time.Time        `json:"createdAt"`
	Status               int        `json:"status"`
}

func (consumption *Consumption) Map(ticket soda.Ticket, user soda_manager.User, parentUser soda_manager.User, device soda_manager.Device, paymentList []*soda.Payment) {
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
		consumption.TypeName = "/"
	}
	consumption.Value = ticket.Value
	for _, payment := range paymentList {
		if payment.ID == ticket.PaymentId {
			consumption.Payment = payment.Name
			break
		}
	}
	if consumption.Payment == "" {
		consumption.Payment = "/"
	}
	consumption.PaymentID = ticket.PaymentId
	consumption.CreatedAt = ticket.CreatedAt
}
