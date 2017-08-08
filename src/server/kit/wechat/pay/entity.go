package pay

/**
统一下单请求
*/
type UnifyOrderRequest struct {
	AppId          string `xml:"appid"`
	Body           string `xml:"body"`
	MchId          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	NotifyUrl      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
	TotalFee       int    `xml:"total_fee"`
	ProductId      string `xml:"product_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	Sign           string `xml:"sign"`
	OpenId         string `xml:"openid"`
	Attach         string `xml:"attach"`
}

/**
统一下单返回(已补全)
*/
/*type UnifyOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeUrl    string `xml:"code_url"`
}*/

/**
异步回调,微信后台请求商户后台(已补全)
*/
/*type NotifyRequest struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	AppId              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           string `xml:"total_fee"`
	SettlementTotalFee string `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          string `xml:"coupon_fee"`
	CouponCount        string `xml:"coupon_count"`
	CouponTypeN        string `xml:"coupon_type_$n"`
	CouponIdN          string `xml:"coupon_id_$n"`
	CouponFeeN         string `xml:"coupon_fee_$n"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
}*/

/**
异步回调,商户返回微信后台数据
*/
type NotifyResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

/**
扫码支付模式一中步骤8中返回微信后台步骤3的请求
*/
/*type NativePayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	PrepayId   string `xml:"prepay_id"`
	ResultCode string `xml:"result_code"`
	ErrCodeDes string `xml:"err_code_des"`
	Sign       string `xml:"sign"`
}*/

/**
扫码支付模式一中步骤3.回调商户设置的支付回调url,微信后台请求商户后台
*/
type NativePayRequest struct {
	AppId       string `xml:"appid"`
	OpenId      string `xml:"openid"`
	MchId       string `xml:"mch_id"`
	IsSubscribe string `xml:"is_subscribe"`
	NonceStr    string `xml:"nonce_str"`
	ProductId   string `xml:"product_id"`
	Sign        string `xml:"sign"`
}

/**
调用查询订单API请求参数
*/
type InitiativeRequest struct {
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	OutTradeNo string `xml:"out_trade_no"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
}

/**
调用查询订单API返回参数(已补全)
*/
/*type InitiativeResponse struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	AppId              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	DeviceInfo         string `xml:"device_info"`
	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	TradeState         string `xml:"trade_state"`
	BankType           string `xml:"bank_type"`
	TotalFee           string `xml:"total_fee"`
	SettlementTotalFee string `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          string `xml:"coupon_fee"`
	CouponCount        string `xml:"coupon_count"`
	CouponTypeN        string `xml:"coupon_type_$n"`
	CouponIdN          string `xml:"coupon_id_$n"`
	CouponFeeN         string `xml:"coupon_fee_$n"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
	TradeStateDesc     string `xml:"trade_state_desc"`
}*/

type RefundRequest struct {
	AppId       string `xml:"appid"`
	MchId       string `xml:"mch_id"`
	NonceStr    string `xml:"nonce_str"`
	OutTradeNo  string `xml:"out_trade_no"`
	OutRefundNo string `xml:"out_refund_no"`
	TotalFee    string `xml:"total_fee"`
	RefundFee   string `xml:"refund_fee"`
	OpUserId    string `xml:"op_user_id"`
	Sign        string `xml:"sign"`
}

/*
type RefundResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	//ErrCode          string `xml:"err_code"`
	//ErrCodeDes       string `xml:"err_code_des"`
	AppId         string `xml:"appid"`
	MchId         string `xml:"mch_id"`
	Nonce         string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	OutRefundNo   string `xml:"out_refund_no"`
	RefundId      string `xml:"refund_id"`
	RefundFee     string `xml:"refund_fee"`
	TotalFee      string `xml:"total_fee"`
	FeeType       string `xml:"fee_type"`
	CashFee       string `xml:"cash_fee"`
}
*/

type BatchPayRequest struct {
	AppId       string `xml:"appid"`
	MchId       string `xml:"mch_id"`
	NonceStr    string `xml:"nonce_str"`
	PartnerTradeNo  string `xml:"partner_trade_no"`
	Openid string `xml:"openid"`
	Amount    string `xml:"amount"`
	Desc   string `xml:"desc"`
	CheckName    string `xml:"check_name"`
	ReUserName    string `xml:"re_user_name"`
	SpbillCreateIp    string `xml:"spbill_create_ip"`
	Sign        string `xml:"sign"`
}

/**
type BatchPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	Nonce      string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode          string `xml:"err_code"` // 错误代码 err_code 否	SYSTEMERROR String(32) 错误码信息
						 // SYSTEMERROR 请使用原单号以及原请求参数重试，否则可能造成重复支付等资金风险
	ErrCodeDes       string `xml:"err_code_des"`

	PartnerTradeNo      string `xml:"partner_trade_no"` // 商户订单号
	PaymentNo string `xml:"payment_no"` // 微信订单号
	PaymentTime   string `xml:"payment_time"` // 微信支付成功时间
}**/
