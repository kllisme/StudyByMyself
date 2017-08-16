package model

func (self *Bill) Mapping(user *User) map[string]interface{} {
	settledAt := ""
	if &self.SettledAt != nil {
		settledAt = ""
	} else {
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
			"type":     self.AccountType,
			"payName":  self.AccountName,
			"name":     self.Account,
			"realName": self.RealName,
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
			"name":    self.Account,
			"realName":self.RealName,
		},
		"totalAmount": self.TotalAmount,
		"status":      self.Status,
		"count":       self.OrderCount,
		"id":          self.ID,
	}
}

