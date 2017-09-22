package soda_manager

func (self *Bill) Mapping(user *User) map[string]interface{} {
	user = user.Mapping()
	isAuto := false
	if self.Mode == 0 { // 0代表自动提现
		isAuto = true
	}
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"updatedAt": self.UpdatedAt,
		"settledAt": self.SettledAt,
		"user": map[string]interface{}{
			"id":          user.ID,
			"name":        user.Name,
			"mobile":      user.Mobile,
			"accountName": user.Account,
			"nickName":    user.Nickname,
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
		"isAuto":      isAuto,
	}
}

func (self *DailyBill) Mapping(user *User) map[string]interface{} {
	user = user.Mapping()
	return map[string]interface{}{
		"createdAt": self.CreatedAt,
		"updatedAt": self.UpdatedAt,
		"settledAt": self.SettledAt,
		"billAt":    self.BillAt,
		"user": map[string]interface{}{
			"id":          user.ID,
			"name":        user.Name,
			"mobile":      user.Mobile,
			"accountName": user.Account,
			"nickName":    user.Nickname,
		},
		"account": map[string]interface{}{
			"type":     self.AccountType,
			"payName":  self.AccountName,
			"name":     self.Account,
			"realName": self.RealName,
		},
		"totalAmount": self.TotalAmount,
		"status":      self.Status,
		"count":       self.OrderCount,
		"id":          self.ID,
	}
}
