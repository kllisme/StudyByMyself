package api

import (
	"strconv"
	"strings"
	"time"

	"bytes"
	"encoding/xml"

	"github.com/bitly/go-simplejson"
	"github.com/fatih/structs"
	"github.com/go-errors/errors"
	"github.com/levigross/grequests"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/entity"
	"maizuo.com/soda/erp/api/src/server/kit/alipay"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"maizuo.com/soda/erp/api/src/server/kit/util"
	"maizuo.com/soda/erp/api/src/server/kit/wechat/pay"
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/service"
)

type BillController struct {
}

// 根据微信支付或者支付宝来获取结算单列表
func (self *BillController) ListByAccountType(ctx *iris.Context) {
	userService := &service.UserService{}
	userCashService := &service.UserCashAccountService{}
	limit, _ := ctx.URLParamInt("limit")      // Default: 10
	offset, _ := ctx.URLParamInt("offset")    //  Default: 0 列表起始位:
	createdAt := ctx.URLParam("createdAt")    // 申请时间
	settledAt := ctx.URLParam("settledAt")    // 结算时间
	keys := ctx.URLParam("keys")              // 运营商名称、帐号名称
	accountType, _ := ctx.URLParamInt("type") // 结算支付类型 1:支付宝 2:微信
	status, _ := ctx.URLParamInt("status")    // 账单状态 1:结算成功 2:等待结算 3:结算中 4:结算失败

	billService := &service.BillService{}

	if accountType == 0 {
		common.Render(ctx, "27080101", nil)
		return
	}
	if limit == 0 {
		limit = 10
	}
	total, err := billService.TotalByAccountType(accountType, status, createdAt, settledAt, keys)
	if err != nil {
		common.Render(ctx, "27080102", err)
		return
	}
	billList, err := billService.ListByAccountType(accountType, status, offset, limit, createdAt, settledAt, keys)
	if err != nil {
		common.Logger.Debugln("billService.ListByAccountType err----------", err)
		common.Render(ctx, "27080103", err)
		return
	}
	objects := make([]interface{}, 0)
	for _, bill := range billList {
		user, err := userService.GetById(bill.UserId)
		if err != nil {
			common.Logger.Debugln("获取账单用户信息失败err----------", err)
			common.Render(ctx, "27080106", err)
			return
		}
		userCashAccount,err := userCashService.BasicByUserId(bill.UserId)
		if err != nil {
			common.Logger.Debugln("获取账单用户信息失败err----------", err)
			common.Render(ctx, "27080106", err)
			return
		}
		objects = append(objects, bill.Mapping(user,userCashAccount))
	}

	common.Render(ctx, "27080100", &entity.PaginationData{
		Pagination: entity.Pagination{Total: total, From: offset, To: offset + limit - 1},
		Objects:    objects,
	})
	return
}

func (self *BillController) BatchPay(ctx *iris.Context) {
	billService := service.BillService{}
	billBatchNoService := &service.BillBatchNoService{}
	params := simplejson.New()
	err := ctx.ReadJSON(params)
	if err != nil {
		common.Logger.Warnln("解析json异常,err : ", err)
		common.Render(ctx, "27080201", "解析json异常")
		return
	}

	billIds, err := params.Get("bills").Array()
	if err != nil {
		common.Logger.Warnln("获取bills异常,err : ", err)
		common.Render(ctx, "27080203", "获取bills异常")
		return
	}
	if len(billIds) == 0 {
		common.Render(ctx, "27080204", "未选择任何账单")
		return
	}
	accountType, err := billService.BillTypeByBatchBill(billIds)
	if err != nil && accountType == -1 {
		common.Logger.Warnln("获取选取的账单结算类型失败,err : ", err)
		common.Render(ctx, "27080202", "获取选取的账单结算类型失败")
		return
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
	if accountType == 1 {
		// 支付宝,生成批次号并拼接支付宝支付的参数
		data, code, err = BatchAlipay(billList)
	} else if accountType == 2 {
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
	err = billService.BatchUpdateSubmitAtById(3, billIds)
	if err != nil {
		common.Logger.Warnln("更新账单为'结算中'失败:", err.Error())
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
	//aliPayBillIds := make([]int, 0)
	billBatchNoList := make([]*model.BillBatchNo, 0)
	batchNum := 0
	batchFee := 0
	var aliPayReqParam map[string]string
	aliPayDetailDataStr := ""

	for _, bill := range billList {
		_remark := bill.CreatedAt.Format("01月02日") + "洗衣结算款"

		aliPayDetailDataStr += bill.BillId + "^" + bill.Account + "^" + bill.RealName +
			"^" + functions.Float64ToString(float64(bill.Amount)/100.00, 2) + "^" + _remark + "|" //组装支付宝支付data_detail
		//aliPayBillIds = append(aliPayBillIds,bill.BillId) //组装需要修改为"结账中"状态的支付宝订单
		batchFee += bill.Amount
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
			common.Logger.Warnln("持久化批次号失败:", err.Error())
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
	common.Logger.Warnln("aliPayDetailDataStr====================", aliPayDetailDataStr)
	alipayKit := alipay.AlipayKit{}
	param["service"] = viper.GetString("pay.aliPay.service.batchTransNotify")
	param["partner"] = viper.GetString("pay.aliPay.partner")
	param["_input_charset"] = viper.GetString("pay.aliPay.inputCharset")
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
	common.Logger.Warnln("batchNum======================", batchNum)
	common.Logger.Warnln("param======================", param)
	return param
}

func (self *BillController) AlipayNotification(ctx *iris.Context) {
	var err error
	billService := &service.BillService{}
	billBatchNoService := &service.BillBatchNoService{}
	billRelService := &service.BillRelService{}
	alipayKit := &alipay.AlipayKit{}
	common.Logger.Warningln("======================支付宝回调开始======================")
	reqMap := make(map[string]string, 0)
	billList := make([]*model.Bill, 0)
	failureList := make([]*model.Bill, 0)
	billRelList := make([]*model.BillRel, 0)
	successedBillIds := make([]string, 0)
	failureBillIds := make([]string, 0)
	successedNotifyDetail := make([]string, 0)
	failNotifyDetail := make([]string, 0)
	billIdSettledAtMap := make(map[string]time.Time)
	successedNum := 0
	failureNum := 0
	reqMap["notify_time"] = ctx.FormValueString("notify_time")
	reqMap["notify_type"] = ctx.FormValueString("notify_type")
	reqMap["notify_id"] = ctx.FormValueString("notify_id")
	reqMap["batch_no"] = ctx.FormValueString("batch_no")
	reqMap["pay_user_id"] = ctx.FormValueString("pay_user_id")
	reqMap["pay_user_name"] = ctx.FormValueString("pay_user_name")
	reqMap["pay_account_no"] = ctx.FormValueString("pay_account_no")
	reqMap["success_details"] = ctx.FormValueString("success_details")
	reqMap["fail_details"] = ctx.FormValueString("fail_details")
	common.Logger.Debugln("signType=============", ctx.FormValueString("sign_type"))
	common.Logger.Warnln("reqMap===========================", reqMap)
	if !alipayKit.VerifySign(reqMap, ctx.FormValueString("sign")) {
		common.Logger.Warningln("回调数据校验失败")
		ctx.Response.SetBodyString("fail")
		return
	}
	common.Logger.Debugln("success")

	//successed status of alipaybill
	if reqMap["success_details"] != "" {
		successedNotifyDetail = strings.Split(reqMap["success_details"], "|")
		if len(successedNotifyDetail) > 0 {
			for _, _detail := range successedNotifyDetail {
				if _detail == "" {
					continue
				}
				_info := strings.Split(_detail, "^")
				if len(_info) > 0 {
					_billId := _info[0] //商家流水号
					//_account := _info[1]    //收款方账号
					//_name := _info[2]       //收款账号姓名
					//_amount := _info[3]     //付款金额
					_flag := _info[4]     //成功或失败标识
					_reason := _info[5]   //成功或失败原因
					_alipayno := _info[6] //支付宝内部流水号
					_time := _info[7]     //完成时间
					_settledAt, _ := time.Parse("20060102150405", _time)
					_bill := &model.Bill{BillId: _billId, SettledAt: _settledAt, Status: 2} //已结账
					_billRel := &model.BillRel{BillId: _billId, BatchNo: reqMap["batch_no"], Type: 1, BillType: 1, IsSuccessed: true, Reason: _reason, OuterNo: _alipayno}
					if _flag == "S" {
						billList = append(billList, _bill)
						billRelList = append(billRelList, _billRel)
						successedBillIds = append(successedBillIds, _billId)
						billIdSettledAtMap[_billId] = _settledAt
						successedNum++
					}
				}
			}
		}
	}
	//failure status of alipaybill
	if reqMap["fail_details"] != "" {
		failNotifyDetail = strings.Split(reqMap["fail_details"], "|")
		if len(failNotifyDetail) > 0 {
			for _, _detail := range failNotifyDetail {
				if _detail == "" {
					continue
				}
				_info := strings.Split(_detail, "^")
				if len(_info) > 0 {
					_billId := _info[0]
					_flag := _info[4]
					_reason := _info[5]
					_alipayno := _info[6]
					_time := _info[7]
					_settledAt, _ := time.Parse("20060102150405", _time)
					_bill := &model.Bill{BillId: _billId, SettledAt: _settledAt, Status: 4} //结账失败
					_billRel := &model.BillRel{BillId: _billId, BatchNo: reqMap["batch_no"], Type: 1, BillType: 1, IsSuccessed: false, Reason: _reason, OuterNo: _alipayno}
					if _flag == "F" {
						failureList = append(failureList, _bill)
						billRelList = append(billRelList, _billRel)
						failureBillIds = append(failureBillIds, _billId)
						billIdSettledAtMap[_billId] = _settledAt
						failureNum++
					}
				}
			}
		}
		billList = append(billList, failureList...)
	}
	common.Logger.Debugln("list==============", billList)
	if len(billList) <= 0 {
		common.Logger.Warningln("返回数据没有账单详情")
		ctx.Response.SetBodyString("fail")
		return
	} else {
		err = billService.Updates(&billList)
		if err != nil {
			//更新支付宝账单结账状态失败
			common.Logger.Warnln("更新支付宝账单结账状态失败,原因", ":", err.Error())
			ctx.Response.SetBodyString("fail")
			return
		}
	}

	//软删除失败订单的批次号
	if len(failureBillIds) > 0 {
		_, err = billBatchNoService.Delete(failureBillIds)
		if err != nil {
			common.Logger.Warnln("01060502", "failureBillIds==", failureBillIds, ":", err.Error())
			ctx.Response.SetBodyString("fail")
			return
		}
	}
	//插入支付宝返回的账单信息
	if len(billRelList) > 0 {
		_, err := billRelService.Create(billRelList...)
		if err != nil {
			common.Logger.Warnln("01060503", ":", err.Error())
			ctx.Response.SetBodyString("fail")
			return
		}
	}
	common.Logger.Warnln("回调成功单数:", successedNum, ",失败单数:", failureNum)
	common.Logger.Warningln("======================支付宝回调结束======================")
	ctx.Response.SetBodyString("success")
}

/**
取消提交支付宝批量付款申请,已取消
*/
func (self *BillController) CancelBatchAliPay(ctx *iris.Context) {
	billService := &service.BillService{}
	billBatchNoService := &service.BillBatchNoService{}
	param := simplejson.New()
	err := ctx.ReadJSON(&param)
	if err != nil {
		common.Logger.Debugln("解析json异常")
		common.Render(ctx, "27080301", "解析json异常")
		return
	}
	billIds, err := param.Get("bills").Array()
	if err != nil {
		common.Logger.Debugln("获取参数bills失败")
		common.Render(ctx, "27080302", "获取参数bills失败")
		return
	}
	if len(billIds) <= 0 {
		common.Logger.Debugln("bills为空列表")
		common.Render(ctx, "27080303", "bills为空列表")
		return
	}
	// 先确定这笔账单的状态是结账中,而且账单ID都存在于批次号表
	for _, billId := range billIds {
		bill, err := billService.BasicByBillId(billId.(string))
		if err != nil {
			common.Logger.Debugln("获取账单详情异常,原因:%v,账单ID:%v", err, billId)
			common.Render(ctx, "27080306", err)
			return
		}
		if bill.AccountType != 1 {
			common.Logger.Debugln("存在支付类型不为支付宝的账单,账单ID:%v", billId)
			common.Render(ctx, "27080307", nil)
			return
		}
		if bill.Status != 3 {
			common.Logger.Debugln("存在状态不为'结算中'的账单,账单ID:%v", billId)
			common.Render(ctx, "27080308", nil)
			return
		}
		billBatchNos, err := billBatchNoService.Baisc(billId)
		if err != nil {
			common.Logger.Debugln("获取账单批次详情异常,原因:%v,账单ID:%v", err, billId)
			common.Render(ctx, "27080309", nil)
			return
		}
		if len(*billBatchNos) <= 0 {
			common.Logger.Debugln("账单无批次详情,账单ID:%v", billId)
			common.Render(ctx, "27080310", nil)
			return
		}
	}

	//两个更新暂时没有保持事务
	err = billService.BatchUpdateStatusById(1, billIds) //将"结算中"的状态改成"已申请"
	if err != nil {
		common.Logger.Debugln("更新账单状态'结算中'为'已申请'失败:", billIds)
		common.Render(ctx, "27080304", "更新账单状态'结算中'为'已申请'失败")
		return
	}
	_, err = billBatchNoService.Delete(billIds)
	if err != nil {
		common.Logger.Debugln("取消批次号绑定失败:", billIds)
		common.Render(ctx, "27080305", "取消批次号绑定失败")
		return
	}
	common.Render(ctx, "27080300", nil)
}

func (self *BillController) WechatPay(ctx *iris.Context) {
	billService := &service.BillService{}
	billRelService := &service.BillRelService{}
	billBatchNoService := &service.BillBatchNoService{}
	common.Logger.Warnln("---------------------微信企业支付开始--------------")
	bill, err := billService.GetFirstWechatBill()
	if err != nil {
		common.Render(ctx, "27080401", err)
		return
	}
	status := 3
	billIds := []interface{}{bill.BillId}
	wechatPayKit := pay.WechatPayKit{}
	nonceStr := wechatPayKit.CreateNonceStr(32)
	batchPayRequest := &pay.BatchPayRequest{
		PartnerTradeNo: bill.BillId,
		MchAppId:       viper.GetString("pay.wechat.mchAppId"),
		MchId:          viper.GetString("pay.wechat.mchId"),
		NonceStr:       nonceStr,
		OpenId:         bill.Account,
		CheckName:      viper.GetString("pay.wechat.checkName"),
		ReUserName:     bill.RealName,
		Amount:         bill.Amount,
		Desc:           "企业付款API测试" + bill.CreatedAt.Local().Format("01月02日") + "结算款",
		SPBillCreateIP: "116.24.64.139",
	}
	billRel := &model.BillRel{BillId: bill.BillId, BatchNo: bill.BillId, BillType: 1, Type: 2} // type=2代表微信,billType=1代表记录来源于bill
	respMap, err := BatchWechatPay(batchPayRequest)
	if err != nil {
		common.Render(ctx, "27080402", err)
		return
	}

	if returnCode, ok := respMap["return_code"]; ok == true && returnCode == "FAIL" {
		status = 4
		billRel.IsSuccessed = false
		billRel.Reason = respMap["return_msg"]
		billRel.ErrCode = "return_code : " + respMap["return_code"]
		common.Logger.Warnln("request wechat transfer pay return_code is fail,err : ", respMap["return_msg"])
		billService.BatchUpdateStatusAndSettleAtById(status,billIds)
		_, err = billRelService.Create(billRel)
		if err != nil {
			common.Render(ctx, "27080407", err)
			return
		}
	} else if ok == true && returnCode == "SUCCESS" {
		if resultCode, ok := respMap["result_code"]; ok == true && resultCode == "FAIL" {
			billRel.IsSuccessed = false
			billRel.ErrCode = respMap["err_code"]
			billRel.Reason = respMap["err_code_des"]
			common.Log(ctx, &common.Result{
				Status:      "UNPROCESSABLE_ENTITY",
				Data:        errors.Errorf("err_code : %v,reason : %v", respMap["err_code"], respMap["err_code_des"]),
				Description: "请求微信企业支付通信成功但业务失败",
				Msg:         "请求微信企业支付通信成功但业务失败",
				Exception:   "",
				Code:        "27080408",
				IsError:     false,
			})
			status = 4
			billService.BatchUpdateStatusAndSettleAtById(status,billIds)
			_, err = billRelService.Create(billRel)
			if err != nil {
				common.Render(ctx, "27080407", err)
				return
			}
		} else {
			billRel.ErrCode = respMap["result_code"]
			if respMap["partner_trade_no"] != bill.BillId {
				status = 3
				billRel.IsSuccessed = false
				billRel.Reason = "商户订单号不一致"
				common.Log(ctx, &common.Result{
					Status:      "UNPROCESSABLE_ENTITY",
					Data:        errors.Errorf("bill_id : %v,reason : %v", bill.BillId, respMap["partner_trade_no"]),
					Description: "商户订单号不一致",
					Msg:         "商户订单号不一致",
					Exception:   "",
					Code:        "27080409",
					IsError:     false,
				})

			}else if respMap["nonce_str"] != nonceStr {
				// 返回的随机串有问题
				status = 3
				billRel.IsSuccessed = false
				billRel.Reason = "随机串校验不通过"
				common.Log(ctx, &common.Result{
					Status:      "UNPROCESSABLE_ENTITY",
					Data:        errors.Errorf("nonceStr : %v,respMap['nonce_str'] : %v", nonceStr, respMap["nonce_str"]),
					Description: "随机串校验不通过",
					Msg:         "随机串校验不通过",
					Exception:   "",
					Code:        "27080410",
					IsError:     false,
				})
			} else {
				billRel.IsSuccessed = true
				billRel.Reason = respMap["return_msg"]
				billRel.OuterNo = respMap["payment_no"]
				status = 2
				common.Logger.Warnln("微信企业支付业务成功,payment_no,", respMap["payment_no"])
				billService.BatchUpdateStatusAndSettleAtById(status,billIds)
				if err != nil {
					common.Render(ctx, "27080407", err)
					return
				}
			}
			billService.BatchUpdateStatusById(status,billIds)
			_, err = billRelService.Create(billRel)
			if err != nil {
				common.Render(ctx, "27080407", err)
				return
			}
		}
	} else {
		billRel.IsSuccessed = false
		billRel.ErrCode = ""
		billRel.Reason = "微信企业支付return_code数据出错"
		status = 4
		common.Log(ctx, &common.Result{
			Status:      "UNPROCESSABLE_ENTITY",
			Data:        respMap,
			Description: "the most serious condition--wechat return an err packet",
			Msg:         "the most serious condition--wechat return an err packet",
			Exception:   "",
			Code:        "27080411",
			IsError:     false,
		})
		billService.BatchUpdateStatusAndSettleAtById(status,billIds)
		_, err = billRelService.Create(billRel)
		if err != nil {
			common.Render(ctx, "27080407", err)
			return
		}
	}
	// 表示业务失败了
	if billRel.IsSuccessed != true {
		//软删除失败订单的批次号
		_, err = billBatchNoService.Delete(bill.BillId)
		if err != nil {
			common.Logger.Debugln("01060502", "failureBillIds==", bill.BillId, ":", err.Error())
			return
		}

	}

	common.Render(ctx, "27080400", nil)
	common.Logger.Warningln("---------------------微信企业支付完成--------------")
	return
}

func BatchWechatPay(batchPayRequest *pay.BatchPayRequest) (map[string]string, error) {
	common.Logger.Warningln("batchPayRequest--------------------->", batchPayRequest)
	m := structs.Map(batchPayRequest)
	delete(m, "sign")
	wechatPayKit := pay.WechatPayKit{}
	batchPayRequest.Sign = wechatPayKit.CreateSign(m)
	requestBytes, err := xml.Marshal(batchPayRequest)
	if err != nil {
		common.Logger.Warnln("微信企业支付请求转XML失败,error=========", err)
		return nil, err
	}
	reqStr := string(requestBytes)
	common.Logger.Warnln("reqStr:", reqStr)
	client, err := wechatPayKit.CreateTLSClient(
		viper.GetString("pay.wechat.tlsFile.cert"),
		viper.GetString("pay.wechat.tlsFile.key"),
		viper.GetString("pay.wechat.tlsFile.root"),
	)
	if err != nil {
		common.Logger.Warnln("微信企业支付请求CreateTLSClient失败,error=========", err)
		return nil, err
	}
	// client.Timeout=xxx
	requestBytes = []byte(reqStr)
	url := viper.GetString("pay.wechat.requestUrl.createTransfers")
	response, err := grequests.Post(url, &grequests.RequestOptions{
		XML: requestBytes,
		Headers: map[string]string{
			"Accept":       "application/xml",
			"Content-Type": "application/xml;charset=utf-8",
		},
		HTTPClient: client,
	})
	common.Logger.Warnln("response：", response.String())
	if err != nil {
		common.Logger.Warnln("微信企业支付请求失败,err:", err.Error(), ",statusCode:", response.StatusCode)
		return nil, err
	}
	if response.StatusCode != 200 {
		common.Logger.Warnln("微信企业支付请求返回错误码,statusCode:", response.StatusCode)
		return nil, errors.New(response.StatusCode)

	}
	respMap, err := util.DecodeXMLToMap(bytes.NewReader(response.Bytes()))
	if err != nil {
		common.Logger.Warnln("解析xml形式编码错误, 原因:", err.Error())
		return nil, err
	}
	common.Logger.Warningln("微信企业支付响应：", respMap)
	return respMap, nil
}
