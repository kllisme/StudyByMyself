package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-errors/errors"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/kit/alipay"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"maizuo.com/soda/erp/api/src/server/kit/wechat/pay"
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/service"
	"maizuo.com/soda/erp/api/src/server/service/permission"
)

type BillController struct {
}

// 根据微信支付或者支付宝来获取结算单列表
func (self *BillController) ListByAccountType(ctx *iris.Context) {
	// 或许可以放到中间件
	userRoleService := &permission.UserRoleRelService{}
	userService := &service.UserService{}
	userId, _ := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	userRoleList, err := userRoleService.GetRoleIDsByUserID(userId)
	if err != nil {
		common.Logger.Debugln("获取当前操作用户角色失败 userId------", userId)
		common.Render(ctx, "27070104", err)
		return
	}
	// 判断是不是财务或者系统管理员,不是财务的不放行
	if functions.FindIndex(userRoleList, 3) == -1 && functions.FindIndex(userRoleList, 5) == -1 {
		common.Logger.Debugln("获取当前操作用户不具有权限 userRoleList-----", userRoleList)
		common.Render(ctx, "27070105", nil)
		return
	}

	limit, _ := ctx.URLParamInt("limit")      // Default: 10
	offset, _ := ctx.URLParamInt("offset")    //  Default: 0 列表起始位:
	createdAt := ctx.URLParam("createdAt")    // 申请时间
	settledAt := ctx.URLParam("settledAt")    // 结算时间
	keys := ctx.URLParam("keys")              // 运营商名称、帐号名称
	accountType, _ := ctx.URLParamInt("type") // 结算支付类型 1:支付宝 2:微信
	status, _ := ctx.URLParamInt("status")    // 账单状态 1:结算成功 2:等待结算 3:结算中 4:结算失败

	billService := &service.BillService{}

	if accountType == -1 {
		common.Render(ctx, "27070101", nil)
		return
	}
	if offset == -1 {
		offset = 0
	}
	if limit == -1 {
		limit = 10
	}
	total, err := billService.TotalByAccountType(accountType, status, createdAt, settledAt, keys)
	if err != nil {
		common.Render(ctx, "27070103", err)
	}
	billList, err := billService.ListByAccountType(accountType, status, offset, limit, createdAt, settledAt, keys)
	if err != nil {
		common.Logger.Debugln("billService.ListByAccountType err----------", err)
		common.Render(ctx, "27070102", err)
		return
	}
	objects := make([]interface{},0)
	for _,bill := range  billList {
		user,err := userService.GetById(bill.ID)
		if err != nil {
			common.Logger.Debugln("获取账单用户信息失败err----------", err)
			common.Render(ctx, "27070106", err)
		}
		objects = append(objects,bill.Mapping(user))
	}

	common.Render(ctx, "27070100", &entity.PaginationData{
		Pagination: entity.Pagination{Total: total, From: offset, To: offset + limit},
		Objects:    objects,
	})
	return
}

func (self *BillController) BatchPay(ctx *iris.Context) {
	billService := service.BillService{}
	billBatchNoService := &service.BillBatchNoService{}
	params := simplejson.New()
	if ctx.ReadJSON(params) != nil {
		common.Logger.Debugln("解析json异常")
		common.Render(ctx, "27080201", "解析json异常")
		return
	}
	// 结算支付类型 1:支付宝 2:微信
	payType, err := params.Get("type").Int()
	if err == nil {
		common.Logger.Debugln("获取结算支付类型异常")
		common.Render(ctx, "27080202", "获取结算支付类型异常")
		return
	}
	bills, err := params.Get("bills").Array()
	if err == nil {
		common.Logger.Debugln("获取bills异常")
		common.Render(ctx, "27080203", "获取bills异常")
		return
	}
	if len(bills) == 0 {
		common.Render(ctx, "27080204", "未选择任何账单")
		return
	}
	billIds := make([]interface{}, 0)
	for _, _param := range bills {
		_map := _param.(map[string]interface{})
		billIds = append(billIds, _map["billId"])
	}
	// 确定选取的是发起结算和结账失败的账单
	statusList := []interface{}{1, 4}
	billList, err := billService.ListByBillIdsAndStatus(billIds, statusList)
	if err != nil {
		common.Render(ctx, "27080205", "获取账单列表失败")
		return
	}
	if len(billList) != len(billIds) {
		common.Render(ctx, "27080206", "所选账单中包含不是发起提现和结算失败的订单")
	}
	//查询账单列表中是否已有批次号的订单(再次确认,这里的订单号只是"已申请"和"结账失败"的)
	batchNoList, _ := billBatchNoService.Baisc(billIds)
	if len(*batchNoList) > 0 {
		//common.Render(ctx,"CODE","所选账单中包含已结账账单，请重新选择")
		common.Render(ctx, "27080207", "所选账单中包含已结账账单，请重新选择")
		return
	}
	code, data := "", make(map[string]string)
	if payType == 1 {
		// 支付宝,生成批次号并拼接支付宝支付的参数
		data, code, err = BatchAlipay(billList)
	} else if payType == 2 {
		// 微信,无需做处理,只需要在下面统一改变bill和daily_bill的状态
		err = nil
	} else {
		common.Render(ctx, "27080208", "错误结算类型")
		return
	}

	// 将支付宝处理后的问题排解
	if err != nil {
		common.Render(ctx, code, err)
	}
	// 不用区分微信还是支付宝的单,统一改变bill和daily_bill的状态
	err = billService.BatchUpdateStatusById(3, billIds)
	if err != nil {
		common.Logger.Debugln("更新账单为'结算中'失败:", err.Error())
		common.Render(ctx, "27080209", err)
		return
	}
	// "日账单结账成功"
	common.Render(ctx, "27080200", data)
}

/*
	map[string]string aliPayReqParam
	string   code
*/
func BatchAlipay(billList []*model.Bill) (map[string]string, string, error) {

	billBatchNoService := &service.BillBatchNoService{}
	aliPayBillIds := make([]int, 0)
	billBatchNoList := make([]*model.BillBatchNo, 0)
	batchNum := 0
	batchFee := 0
	var aliPayReqParam map[string]string
	aliPayDetailDataStr := ""

	for _, bill := range billList {
		_remark := bill.CreatedAt.Format("01月02日") + "洗衣结算款"

		aliPayDetailDataStr += strconv.Itoa(bill.ID) + "^" + bill.Account + "^" + bill.RealName +
			"^" + functions.Float64ToString(float64(bill.TotalAmount)/100.00, 2) + "^" + _remark + "|" //组装支付宝支付data_detail
		aliPayBillIds = append(aliPayBillIds, bill.ID) //组装需要修改为"结账中"状态的支付宝订单
		batchFee += bill.TotalAmount
		batchNum++ //计算批量结算请求中支付宝结算的日订单数,不可超过1000
	}

	//生成支付宝请求参数并存储账单对应的批次号
	if batchNum > 0 && batchNum <= 1000 && aliPayDetailDataStr != "" {
		aliPayReqParam = GenerateBatchAliPay(batchNum, batchFee, aliPayDetailDataStr)
		if aliPayReqParam["batch_no"] == "" {
			common.Logger.Debugln("生成批次号失败")
			return nil, "27080210", errors.New("支付宝结算更新状态失败")
		}
		//create bill_batch_no
		for _, _bill := range billList {
			_billBatchNo := &model.BillBatchNo{
				BillId:   _bill.BillId,
				BatchNo:  aliPayReqParam["batch_no"],
				BillType: 1, // 1为bill
			}
			billBatchNoList = append(billBatchNoList, _billBatchNo)
		}
		if len(billBatchNoList) <= 0 {
			common.Logger.Debugln("生成批次号信息失败")
			return nil, "27080211", errors.New("支付宝结算更新状态失败")
		}
		_, err := billBatchNoService.BatchCreate(&billBatchNoList)
		if err != nil {
			common.Logger.Debugln("持久化批次号失败:", err.Error())
			return nil, "27080212", errors.New("支付宝结算更新状态失败")
		}

	} else if (batchNum <= 0 || batchNum > 1000) && aliPayDetailDataStr != "" {
		common.Logger.Debugln("所选支付宝账单超出批次结算最大值1000", ":", batchNum)
		return nil, "27080213", errors.New("所选支付宝账单超出批次结算最大值1000")
	}
	return aliPayReqParam, "", nil
}

/*
	支付宝批量支付接口
*/
func GenerateBatchAliPay(batchNum int, batchFee int, aliPayDetailDataStr string) map[string]string {

	param := make(map[string]string, 0)
	if batchNum <= 0 || batchNum > 1000 || aliPayDetailDataStr == "" {
		return param
	}
	if strings.HasSuffix(aliPayDetailDataStr, "|") {
		aliPayDetailDataStr = aliPayDetailDataStr[:len(aliPayDetailDataStr)-1]
	}
	common.Logger.Debugln("aliPayDetailDataStr====================", aliPayDetailDataStr)
	alipayKit := alipay.AlipayKit{}
	param["service"] = "batch_trans_notify"
	param["partner"] = viper.GetString("pay.aliPay.id")
	param["_input_charset"] = "utf-8"
	param["notify_url"] = viper.GetString("pay.aliPay.notifyUrl")
	param["account_name"] = viper.GetString("pay.aliPay.accountName")
	param["detail_data"] = aliPayDetailDataStr
	param["batch_no"] = time.Now().Local().Format("20060102150405")
	param["batch_num"] = strconv.Itoa(batchNum)
	param["batch_fee"] = functions.Float64ToString(float64(batchFee)/100.00, 2)
	param["email"] = viper.GetString("pay.aliPay.email")
	param["pay_date"] = time.Now().Local().Format("20060102")
	param["sign"] = alipayKit.CreateSign(param)
	param["sign_type"] = viper.GetString("pay.aliPay.signType")
	param["request_url"] = viper.GetString("pay.aliPay.requestUrl")
	common.Logger.Debugln("batchNum======================", batchNum)
	common.Logger.Debugln("param======================", param)
	return param
}

//func (self *DailyBillController) AlipayNotification(ctx *iris.Context) {
//	var err error
//	dailyBillService := &service.DailyBillService{}
//	billBatchNoService := &service.BillBatchNoService{}
//	billRelService := &service.BillRelService{}
//	common.Logger.Warningln("======================支付宝回调开始======================")
//	reqMap := make(map[string]string, 0)
//	billList := make([]*model.DailyBill, 0)
//	failureList := make([]*model.DailyBill, 0)
//	billRelList := make([]*model.BillRel, 0)
//	successedBillIds := make([]int, 0)
//	failureBillIds := make([]int, 0)
//	successedNotifyDetail := make([]string, 0)
//	failNotifyDetail := make([]string, 0)
//	mnznBillList := make([]map[string]interface{}, 0)
//	billIdSettledAtMap := make(map[string]string)
//	successedNum := 0
//	failureNum := 0
//	reqMap["notify_time"] = ctx.FormValueString("notify_time")
//	reqMap["notify_type"] = ctx.FormValueString("notify_type")
//	reqMap["notify_id"] = ctx.FormValueString("notify_id")
//	reqMap["batch_no"] = ctx.FormValueString("batch_no")
//	reqMap["pay_user_id"] = ctx.FormValueString("pay_user_id")
//	reqMap["pay_user_name"] = ctx.FormValueString("pay_user_name")
//	reqMap["pay_account_no"] = ctx.FormValueString("pay_account_no")
//	reqMap["success_details"] = ctx.FormValueString("success_details")
//	reqMap["fail_details"] = ctx.FormValueString("fail_details")
//	common.Logger.Debugln("signType=============", ctx.FormValueString("sign_type"))
//	common.Logger.Debugln("reqMap===========================", reqMap)
//	if !alipay.VerifySign(reqMap, ctx.FormValueString("sign")) {
//		common.Logger.Warningln("回调数据校验失败")
//		ctx.Response.SetBodyString("fail")
//		return
//	}
//	common.Logger.Debugln("success")
//
//	//successed status of alipaybill
//	if reqMap["success_details"] != "" {
//		successedNotifyDetail = strings.Split(reqMap["success_details"], "|")
//		if len(successedNotifyDetail) > 0 {
//			for _, _detail := range successedNotifyDetail {
//				if _detail == "" {
//					continue
//				}
//				_info := strings.Split(_detail, "^")
//				if len(_info) > 0 {
//					_id := _info[0] //商家流水号
//					//_account := _info[1]    //收款方账号
//					//_name := _info[2]       //收款账号姓名
//					//_amount := _info[3]     //付款金额
//					_flag := _info[4]     //成功或失败标识
//					_reason := _info[5]   //成功或失败原因
//					_alipayno := _info[6] //支付宝内部流水号
//					_time := _info[7]     //完成时间
//					insertTime, _ := time.Parse("20060102150405", _time)
//					_settledAt := insertTime.Format("2006-01-02 15:04:05")
//					_bill := &model.DailyBill{Model: model.Model{Id: functions.StringToInt(_id)}, SettledAt: _settledAt, Status: 2} //已结账
//					_billRel := &model.BillRel{BillId: _id, BatchNo: reqMap["batch_no"], Type: 1, IsSuccessed: true, Reason: _reason, OuterNo: _alipayno}
//					if _flag == "S" {
//						billList = append(billList, _bill)
//						billRelList = append(billRelList, _billRel)
//						successedBillIds = append(successedBillIds, functions.StringToInt(_id))
//						billIdSettledAtMap[_id] = _settledAt
//						successedNum++
//					}
//				}
//			}
//		}
//	}
//	//failure status of alipaybill
//	if reqMap["fail_details"] != "" {
//		failNotifyDetail = strings.Split(reqMap["fail_details"], "|")
//		if len(failNotifyDetail) > 0 {
//			for _, _detail := range failNotifyDetail {
//				if _detail == "" {
//					continue
//				}
//				_info := strings.Split(_detail, "^")
//				if len(_info) > 0 {
//					_id := _info[0]
//					_flag := _info[4]
//					_reason := _info[5]
//					_alipayno := _info[6]
//					_time := _info[7]
//					insertTime, _ := time.Parse("20060102150405", _time)
//					_settledAt := insertTime.Format("2006-01-02 15:04:05")
//					_bill := &model.DailyBill{Model: model.Model{Id: functions.StringToInt(_id)}, SettledAt: _settledAt, Status: 4} //结账失败
//					_billRel := &model.BillRel{BillId: _id, BatchNo: reqMap["batch_no"], Type: 1, IsSuccessed: false, Reason: _reason, OuterNo: _alipayno}
//					if _flag == "F" {
//						failureList = append(failureList, _bill)
//						billRelList = append(billRelList, _billRel)
//						failureBillIds = append(failureBillIds, functions.StringToInt(_id))
//						billIdSettledAtMap[_id] = _settledAt
//						failureNum++
//					}
//				}
//			}
//		}
//		billList = append(billList, failureList...)
//	}
//	common.Logger.Debugln("list==============", billList)
//	if len(billList) <= 0 {
//		common.Logger.Warningln("返回数据没有账单详情")
//		ctx.Response.SetBodyString("fail")
//		return
//	}
//
//	//用于更新旧系统数据
//	if len(successedBillIds) > 0 {
//		mnznSuccessedBillList, err := BasicMnznBillLists(true, &billIdSettledAtMap, successedBillIds...)
//		if err != nil {
//			common.Logger.Debugln(daily_bill_msg["01060504"], ":", err.Error())
//			result := &entity.Result{"01060504", err.Error(), daily_bill_msg["01060504"]}
//			common.Log(ctx, result)
//			ctx.Response.SetBodyString("fail")
//			return
//		}
//		mnznBillList = append(mnznBillList, mnznSuccessedBillList...)
//	}
//	if len(failureBillIds) > 0 {
//		mnznFailureBillList, err := BasicMnznBillLists(false, &billIdSettledAtMap, failureBillIds...)
//		if err != nil {
//			common.Logger.Debugln(daily_bill_msg["01060505"], ":", err.Error())
//			result := &entity.Result{"01060505", err.Error(), daily_bill_msg["01060505"]}
//			common.Log(ctx, result)
//			ctx.Response.SetBodyString("fail")
//			return
//		}
//		mnznBillList = append(mnznBillList, mnznFailureBillList...)
//	}
//
//	if len(mnznBillList) > 0 {
//		_, err = dailyBillService.Updates(&billList, mnznBillList)
//		if err != nil {
//			//更新支付宝账单结账状态失败
//			common.Logger.Debugln(daily_bill_msg["01060501"], ":", err.Error())
//			result := &entity.Result{"01060501", err.Error(), daily_bill_msg["01060501"]}
//			common.Log(ctx, result)
//			ctx.Response.SetBodyString("fail")
//			return
//		}
//	}
//	//软删除失败订单的批次号
//	if len(failureBillIds) > 0 {
//		_, err = billBatchNoService.Delete(failureBillIds)
//		if err != nil {
//			common.Logger.Debugln(daily_bill_msg["01060502"], "failureBillIds==", failureBillIds, ":", err.Error())
//			result := &entity.Result{"01060502", err.Error(), daily_bill_msg["01060502"]}
//			common.Log(ctx, result)
//			ctx.Response.SetBodyString("fail")
//			return
//		}
//	}
//	//插入支付宝返回的账单信息
//	if len(billRelList) > 0 {
//		_, err := billRelService.Create(billRelList...)
//		if err != nil {
//			common.Logger.Debugln(daily_bill_msg["01060503"], ":", err.Error())
//			result := &entity.Result{"01060503", err.Error(), daily_bill_msg["01060503"]}
//			common.Log(ctx, result)
//			ctx.Response.SetBodyString("fail")
//			return
//		}
//	}
//	common.Logger.Warningln("======================支付宝回调结束======================")
//	common.Log(ctx, nil)
//	ctx.Response.SetBodyString("success")
//}

func (self *BillController) WechatPay(bills interface{}) {
	billService := &service.BillService{}
	bill, err := billService.GetFirstWechatBill()
	if err != nil {
		common.Logger.Debugln("获取微信账单失败")
	}
	status := 3
	billIds := []interface{}{bill.BillId}
	wechatPayKit := pay.WechatPayKit{}
	batchPayRequest := &pay.BatchPayRequest{}
	batchPayRequest.PartnerTradeNo = bill.BillId
	batchPayRequest.Desc = bill.CreatedAt.Local().Format("01月02日") + "结算款"
	batchPayRequest.Amount = strconv.Itoa(bill.Amount)
	batchPayRequest.ReUserName = bill.AccountName
	batchPayRequest.Openid = bill.Account
	respMap,err := wechatPayKit.BatchPay(batchPayRequest)
	if err != nil {
		common.Logger.Debugln("微信企业支付失败")
		status = 4
	}
	if respMap["return_code"] == "FAIL" {
		common.Logger.Debugln("请求微信企业支付成功但通信失败,原因:",respMap["return_msg"])
		status = 4
	}else{
		if respMap["result_code"] == "FAIL"  {
			if respMap["err_code"] == "SYSTEMERROR" {
				// 系统错误，请重试	请使用原单号以及原请求参数重试，否则可能造成重复支付等资金风险
			}else {
				common.Logger.Debugln("请求微信企业支付通信成功但业务失败,错误码",
					respMap["err_code"],",原因:",respMap["err_code_des"])
				status = 4
			}
		}else{
			status = 3
		}
	}
	billService.BatchUpdateStatusById(status,billIds)
}
