package soda

import mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"

func (self *Ticket) Mapping(device *mngModel.Device, dailyBill *mngModel.DailyBill) map[string]interface{} {
	settledAt := ""
	if &dailyBill.SettledAt != nil {
		settledAt = ""
	} else {
		settledAt = dailyBill.SettledAt.Format("2006-01-02T15:04:05+00:00")
	}
	if self.Status == 4 || self.PaymentId == 4 {
		self.Settle = false
	} else {
		self.Settle = true
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
