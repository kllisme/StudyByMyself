package model

type DailyOperate struct {
	Model
	Date                 string `json:"date"`
	TotalRecharge        int    `json:"totalRecharge"`
	TotalConsume         int    `json:"totalConsume"`
	TotalWechatConsume   int    `json:"totalWechatConsume"`
	TotalAlipayConsume   int    `json:"totalAlipayConsume"`
	TotalWalletConsume   int    `json:"totalWalletConsume"`
	TotalChipcardConsume int    `json:"totalChipcardConsume"`
	TotalDevice          int    `json:"totalDevice"`
	TotalSoldDevice      int    `json:"totalSoldDevice"`
	TotalUnusedDevice    int    `json:"totalUnusedDevice"`
	TotalUser            int    `json:"totalUser"`
	TotalNewUser         int    `json:"totalNewUser"`
	TotalWechatRecharge  int    `json:"totalWechatRecharge"`
	TotalAlipayRecharge  int    `json:"totalAlipayRecharge"`
	TotalRechargeUser    int    `json:"totalRechargeUser"`
	TotalConsumeUser     int    `json:"totalConsumeUser"`
	CreatedTimestamp     int    `json:"createdTimestamp"`
}

func (DailyOperate) TableName() string {
	return "daily_operate"
}
