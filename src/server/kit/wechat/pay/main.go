package pay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/spf13/viper"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
	"maizuo.com/soda/erp-api/src/server/common"
	"maizuo.com/soda/erp-api/src/server/kit/functions"
)

type WechatPayKit struct {
}

/**
微信支付计算签名的函数
*/
func (self *WechatPayKit) CreateSign(mReq map[string]interface{}) (string) {
	apiKey := viper.GetString("resource.pay.wechat.api-key")
	//apiKey := os.Getenv("APIKEY")
	common.Logger.Debugln("微信支付签名计算, API KEY:", apiKey)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	if apiKey != "" {
		signStrings = signStrings + "key=" + apiKey
	}
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

/**
签名校验方法
*/
func (self *WechatPayKit) VerifySign(data map[string]interface{}, sign string) bool {
	_sign := self.CreateSign(data)
	common.Logger.Debugln("计算出来的sign: %v", _sign)
	common.Logger.Debugln("微信通知sign: %v", sign)
	if sign == _sign {
		common.Logger.Debugln("签名校验通过!")
		return true
	}
	common.Logger.Debugln("签名校验失败!")
	return false
}

/**
创建随机字符串
*/
func (self *WechatPayKit) CreateNonceStr(_len int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < _len; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	text := string(result)
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

/**
创建二维码支付地址
*/
func (self *WechatPayKit) CreateNativePayURL(tradeNo string) string {
	nonceStr := self.CreateNonceStr(32)
	appId := viper.GetString("resource.pay.wechat.app-id")
	mchId := viper.GetString("resource.pay.wechat.mch-id")    //商户号
	//appId := os.Getenv("APPID")
	//mchId := os.Getenv("MCHID")
	timeStamp := time.Now().In(time.FixedZone("Asia/Shanghai", 8*60*60)).Unix() //.Format("20060102150405")
	m := make(map[string]interface{}, 0)
	m["appid"] = appId
	m["mch_id"] = mchId
	m["product_id"] = tradeNo
	m["time_stamp"] = timeStamp
	m["nonce_str"] = nonceStr
	sign := self.CreateSign(m)
	return "weixin://wxpay/bizpayurl?sign=" + sign +
		"&appid=" + url.QueryEscape(appId) +
		"&mch_id=" + url.QueryEscape(mchId) +
		"&product_id=" + url.QueryEscape(tradeNo) +
		"&time_stamp=" + url.QueryEscape(functions.Int64ToString(timeStamp)) +
		"&nonce_str=" + url.QueryEscape(nonceStr)
}

/**
统一下单接口
*/
func (self *WechatPayKit) CreateUnifiedOrder(unifyOrderRequest *UnifyOrderRequest) (string, error) {
	appId := viper.GetString("resource.pay.wechat.app-id")
	//appId := os.Getenv("APPID")
	notifyUrl := viper.GetString("resource.pay.wechat.notifty-url")
	spbillCreateIp := viper.GetString("resource.pay.wechat.ip")
	unifiedorderUrl := viper.GetString("resource.pay.wechat.unifiedorder-url")
	mchId := viper.GetString("resource.pay.wechat.mch-id")
	unifyOrderRequest.AppId = appId
	unifyOrderRequest.NotifyUrl = notifyUrl
	unifyOrderRequest.SpbillCreateIp = spbillCreateIp
	unifyOrderRequest.MchId = mchId

	m := make(map[string]interface{}, 0)
	m["appid"] = unifyOrderRequest.AppId
	m["body"] = unifyOrderRequest.Body
	m["mch_id"] = unifyOrderRequest.MchId
	m["notify_url"] = unifyOrderRequest.NotifyUrl
	m["trade_type"] = unifyOrderRequest.TradeType
	m["spbill_create_ip"] = unifyOrderRequest.SpbillCreateIp
	m["total_fee"] = unifyOrderRequest.TotalFee
	m["out_trade_no"] = unifyOrderRequest.OutTradeNo
	m["nonce_str"] = unifyOrderRequest.NonceStr
	m["product_id"] = unifyOrderRequest.ProductId
	m["openid"] = unifyOrderRequest.OpenId
	common.Logger.Debugln("openid=======", unifyOrderRequest.OpenId)
	unifyOrderRequest.Sign = self.CreateSign(m)
	requestBytes, err := xml.Marshal(unifyOrderRequest)
	if err != nil {
		common.Logger.Warningln("以xml形式编码发送错误, 原因:", err.Error())
		return "", err
	}
	reqStr := string(requestBytes)
	reqStr = strings.Replace(reqStr, "UnifyOrderRequest", "xml", -1)

	common.Logger.Debugln("统一下单接口请求参数:", reqStr)

	requestBytes = []byte(reqStr)
	response, err := grequests.Post(unifiedorderUrl, &grequests.RequestOptions{
		XML: requestBytes,
		Headers: map[string]string{
			"Accept":       "application/xml",
			"Content-Type": "application/xml;charset=utf-8",
		},
	})
	if err != nil {
		common.Logger.Warningln("请求微信支付统一下单接口发送错误, 原因:", err.Error())
		return "", err
	}
	common.Logger.Infoln(response.String())

	xmlResp := UnifyOrderResponse{}
	err = xml.Unmarshal(response.Bytes(), &xmlResp)
	if err != nil {
		common.Logger.Warningln("解析xml形式编码错误, 原因:", err.Error())
		return "", err
	}
	if xmlResp.ReturnCode == "FAIL" {
		common.Logger.Warningln("微信支付统一下单不成功，原因:", xmlResp.ReturnCode)
		return "", err
	}

	//拿prepayId要判断returncode和resultcode还有prepayid
	common.Logger.Warningln("微信支付统一下单成功，预支付单号:", xmlResp.PrepayId)
	return xmlResp.PrepayId, nil
}

/**
查询订单
*/
func (self *WechatPayKit) CheckTrade(tradeNo string) (*InitiativeResponse, error) {
	appId := viper.GetString("resource.pay.wechat.app-id")
	mchId := viper.GetString("resource.pay.wechat.mch-id")
	initiativeUrl := viper.GetString("resource.pay.wechat.initiative-url")
	xmlResp := InitiativeResponse{}
	nonceStr := self.CreateNonceStr(32)
	m := make(map[string]interface{}, 0)

	m["appid"] = appId
	m["mch_id"] = mchId
	m["out_trade_no"] = tradeNo
	m["nonce_str"] = nonceStr
	sign := self.CreateSign(m)
	initiativeRequest := InitiativeRequest{
		AppId:      appId,
		MchId:      mchId,
		OutTradeNo: tradeNo,
		NonceStr:   nonceStr,
		Sign:       sign,
	}
	requestBytes, err := xml.Marshal(initiativeRequest)
	if err != nil {
		common.Logger.Warningln("以xml形式编码发送错误, 原因:", err.Error())
		return nil, err
	}
	reqStr := string(requestBytes)
	reqStr = strings.Replace(reqStr, "InitiativeRequest", "xml", -1)
	response, err := grequests.Post(initiativeUrl, &grequests.RequestOptions{
		XML: requestBytes,
		Headers: map[string]string{
			"Accept":       "application/xml",
			"Content-Type": "application/xml;charset=utf-8",
		},
	})
	if err != nil {
		common.Logger.Warningln("请求微信查询订单接口发送错误, 原因:", err.Error())
		return nil, err
	}

	err = xml.Unmarshal(response.Bytes(), &xmlResp)
	if err != nil {
		common.Logger.Warningln("解析xml形式编码错误, 原因:", err.Error())
		return nil, err
	}

	return &xmlResp, nil

}

func (self *WechatPayKit) Refund(refundRequest *RefundRequest) (*RefundResponse, error) {
	wechatPayKit := &WechatPayKit{}
	refundUrl := viper.GetString("resource.pay.wechat.refund-url")
	mchId := viper.GetString("resource.pay.wechat.mch-id")
	appId := viper.GetString("resource.pay.wechat.app-id")
	m := make(map[string]interface{}, 0)
	/*refundRequest := &pay.RefundRequest{
		AppId: appId,
		MchId: mchId,
		NonceStr: wechatPayKit.CreateNonceStr(32),
		OutTradeNo: trade.TradeId,
		OutRefundNo: order.GenerateIdByMobile(trade.Mobile),
		TotalFee: trade.Value,
		RefundFee: trade.Value,
		OpUserId: mchId,
	}*/
	m["appid"] = appId
	m["mch_id"] = mchId
	m["nonce_str"] = refundRequest.NonceStr
	m["out_trade_no"] = refundRequest.OutTradeNo
	m["out_refund_no"] = refundRequest.OutRefundNo
	m["total_fee"] = refundRequest.TotalFee
	m["refund_fee"] = refundRequest.RefundFee
	m["op_user_id"] = mchId
	refundRequest.Sign = wechatPayKit.CreateSign(m)

	requestBytes, err := xml.Marshal(m)
	if err != nil {
		common.Logger.Warningln("以xml形式编码发送错误, 原因:", err.Error())
		return nil, err
	}
	reqStr := string(requestBytes)
	reqStr = strings.Replace(reqStr, "RefundRequest", "xml", -1)

	requestBytes = []byte(reqStr)
	response, err := grequests.Post(refundUrl, &grequests.RequestOptions{
		XML: requestBytes,
		Headers: map[string]string{
			"Accept":       "application/xml",
			"Content-Type": "application/xml;charset=utf-8",
		},
	})
	if err != nil {
		common.Logger.Warningln("请求微信申请退款接口发送错误, 原因:", err.Error())
		return nil, err
	}
	common.Logger.Infoln(response.String())

	xmlResp := &RefundResponse{}
	err = xml.Unmarshal(response.Bytes(), &xmlResp)
	if err != nil {
		common.Logger.Warningln("解析xml形式编码错误, 原因:", err.Error())
		return nil, err
	}
	return xmlResp, nil
}

func (self *WechatPayKit) UnifiedOrderResponse(nativePayResponse *NativePayResponse) string {
	m := make(map[string]interface{}, 0)
	m["result_code"] = nativePayResponse.ResultCode
	m["return_code"] = nativePayResponse.ReturnCode
	m["return_msg"] = nativePayResponse.ReturnMsg
	m["appid"] = nativePayResponse.AppId
	m["mch_id"] = nativePayResponse.MchId
	m["nonce_str"] = nativePayResponse.NonceStr
	m["prepay_id"] = nativePayResponse.PrepayId
	m["err_code_des"] = nativePayResponse.ErrCodeDes
	nativePayResponse.Sign = self.CreateSign(m)
	respBytes, _ := xml.Marshal(nativePayResponse)
	xmlResult := strings.Replace(string(respBytes), "NativePayResponse", "xml", -1)
	fmt.Println(xmlResult)
	return xmlResult
}
