package model

type UserCashAccount struct {
	Model
	UserId       int    `json:"userId"`
	Type         int    `json:"type"`
	RealName     string `json:"realName"`
	HeadBankName string `json:"headBankName"`
	BankName     string `json:"bankName"`
	Account      string `json:"account"`
	Mobile       string `json:"mobile"`
	CityId       int    `json:"cityId"`
	ProvinceId   int    `json:"provinceId"`
}

func (UserCashAccount) TableName() string {
	return "user_cash_account"
}
