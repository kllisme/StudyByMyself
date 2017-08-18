package soda

import "maizuo.com/soda/erp/api/src/server/model"

func (self *Ticket) Mapping(device *model.Device, dailyBill *model.DailyBill) map[string]interface{} {
	settledAt := ""
	if &dailyBill.SettledAt != nil {
		settledAt = ""
	} else {
		settledAt = dailyBill.SettledAt.Format("2006-01-02T15:04:05+00:00")
	}
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"settledAt": settledAt,
		"user": map[string]interface{}{
			"id":     self.UserId,
			"mobile": self.Mobile,
		},
		"device": map[string]interface{}{
			"serial":  device.SerialNumber,
			"address": device.Address,
		},
		"pay": map[string]interface{}{
			"type":   self.PaymentId,
			"amount": self.Value,
		},
		"type":       self.DeviceMode,
		"id":         self.ID,
		"status":     self.Status,
		"hasSettled": self.Settle,
	}
}
