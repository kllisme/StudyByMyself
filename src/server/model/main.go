package model

func (self *Bill) Mapping(user *User) map[string]interface{} {
	payName := ""
	switch self.AccountType {
	case 1:
		payName = "支付宝"
	case 2:
		payName = "微信"

	}
	settledAt := ""
	if &self.SettledAt != nil {
		settledAt = ""
	}else{
		settledAt = self.SettledAt.Format("2006-01-02T15:04:05+00:00")
	}
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"updatedAt": self.UpdatedAt,
		"settledAt": settledAt,
		"user": map[string]interface{}{
			"id":     user.ID,
			"name":   user.Name,
			"mobile": user.Mobile,
		},
		"account": map[string]interface{}{
			"type":    self.AccountType,
			"payName": payName,
			"name":    self.AccountName,
		},
		"totalAmount": self.TotalAmount,
		"amount":      self.Amount,
		"cast":        self.Cast,
		"rate":        self.Rate / 100,
		"status":      self.Status,
		"count":       self.Count,
		"id":          self.BillId,
	}
}

func (self *DailyBill) Mapping(user *User) map[string]interface{} {
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"updatedAt": self.UpdatedAt,
		"settledAt": self.SettledAt,
		"billAt":    self.BillAt,
		"user": map[string]interface{}{
			"id":     user.ID,
			"name":   user.Name,
			"mobile": user.Mobile,
		},
		"account": map[string]interface{}{
			"type":    self.AccountType,
			"payName": self.AccountName,
			"name":    self.RealName,
		},
		"totalAmount": self.TotalAmount,
		"status":      self.Status,
		"count":       self.OrderCount,
		"id":          self.ID,
	}
}
//func (self *DailyBill) MappingDetails(device *Device, ticket *soda.Ticket) map[string]interface{} {
//	return map[string]interface{}{
//		"createdAt": self.CreatedAt,
//		"updatedAt": self.UpdatedAt,
//		"settledAt": self.SettledAt,
//		"billAt":    self.BillAt,
//		"user": map[string]interface{}{
//			"id":     ticket.UserId,
//			"mobile": ticket.Mobile,
//		},
//		"device": {
//			"serial":  device.SerialNumber,
//			"address": device.Address,
//		},
//		"pay": {
//			"type":   ticket.PaymentId,
//			"amount": self.TotalAmount,
//		},
//		"totalAmount": self.TotalAmount,
//		"status":      ticket.Status,
//		"hasSettled":  self.SettledAt,
//		"id":          self.ID,
//		"type":        ticket.DeviceMode,
//	}
//}
