package soda

import "maizuo.com/soda/erp/api/src/server/model"

func (self *Ticket) Mapping(device *model.Device, dailyBill *model.DailyBill) map[string]interface{} {
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"settledAt": dailyBill.SettledAt,
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
