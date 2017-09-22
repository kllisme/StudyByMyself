package finance

import (
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/excel"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
)

type BillReportController struct {
}

/**获取结算报表详情 */
func (self *BillReportController) DetailsOfReport(ctx *iris.Context) {
	dailyOperateService := &mngService.DailyOperateService{}
	billService := &mngService.BillService{}
	type Response struct {
		Date   time.Time              `json:"date"`
		Wechat map[string]interface{} `json:"wechat"`
		Alipay map[string]interface{} `json:"alipay"`
	}
	startAt := ctx.URLParam("startAt") // 开始时间
	endAt := ctx.URLParam("endAt")     // 结束时间
	if startAt == "" || endAt == "" {
		common.Render(ctx, "27100101", nil)
		return
	}
	start, _ := time.Parse("2006-01-02", startAt[:10])
	start = start.Add(-8 * time.Hour) // 减去八个小时,否则转换的时间戳是北京时间早上8点的时间戳,从而出现检索不到startAt对应那天的记录

	end, _ := time.Parse("2006-01-02", endAt[:10])

	dailyOperateList, err := dailyOperateService.ListByPeriod(start.Local(), end.Local()) // 获取平台收入的记录的map
	if err != nil {
		common.Logger.Warnln("ListByPeriod err ", err)
		common.Render(ctx, "27100102", err)
		return
	}

	alipayType, wechatType := 1, 2
	alipayMap, err := billService.ReportMapByPeriodAndAccountType(startAt[:10], endAt[:10], alipayType)
	if err != nil {
		common.Logger.Warnln("BasicByAccountTypeAndDatetime ------------", err)
		common.Render(ctx, "27100102", err)
		return
	}

	wechatMap, err := billService.ReportMapByPeriodAndAccountType(startAt[:10], endAt[:10], wechatType)
	if err != nil {
		common.Logger.Warnln("BasicByAccountTypeAndDatetime ------------", err)
		common.Render(ctx, "27100102", err)
		return
	}
	responseList := []Response{}
	getValOrDefaultVal := func(m map[string]interface{}, key string) int {
		if value, ok := m[key]; ok != false {
			return value.(int)
		} else {
			return 0
		}
	}
	for _, value := range *dailyOperateList {
		response := Response{}
		// 2016-12-11T12:33:49+08:00
		date, _ := time.Parse("2006-01-02", value.Date)
		response.Date = date.Local()
		wMap := make(map[string]interface{})
		wMap["totalAmount"] = value.TotalWechatConsume + value.TotalWechatRecharge
		wMap["settlement"] = map[string]interface{}{
			"totalAmount": getValOrDefaultVal((*wechatMap)[value.Date], "totalAmount"),
			"cast":        getValOrDefaultVal((*wechatMap)[value.Date], "cast"),
		}
		response.Wechat = wMap

		aMap := make(map[string]interface{})
		aMap["totalAmount"] = value.TotalAlipayConsume + value.TotalAlipayRecharge
		aMap["settlement"] = map[string]interface{}{
			"totalAmount": getValOrDefaultVal((*alipayMap)[value.Date], "totalAmount"),
			"cast":        getValOrDefaultVal((*alipayMap)[value.Date], "cast"),
		}
		response.Alipay = aMap
		responseList = append(responseList, response)
	}
	common.Render(ctx, "27100100", map[string]interface{}{
		"objects": responseList,
	})
	return
}

/** 导出结算报表详情 */
func (self *BillReportController) Export(ctx *iris.Context) {
	dailyOperateService := &mngService.DailyOperateService{}
	billService := &mngService.BillService{}
	type Response struct {
		Date   string                 `json:"date"`
		Wechat map[string]interface{} `json:"wechat"`
		Alipay map[string]interface{} `json:"alipay"`
	}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27100201", err) // json解析失败
	}
	startAt := params.Get("startAt").MustString() // 开始时间
	endAt := params.Get("endAt").MustString()     // 结束时间

	if startAt == "" || endAt == "" {
		common.Render(ctx, "27100202", nil)
		return
	}
	start, _ := time.Parse("2006-01-02", startAt[:10]) // 只取日期,除去时间可能对时间戳造成的影响
	start = start.Add(-8 * time.Hour)                  // 减去八个小时,否则转换的时间戳是北京时间早上8点的时间戳,从而出现检索不到startAt对应那天的记录

	end, _ := time.Parse("2006-01-02", endAt[:10])

	dailyOperateList, err := dailyOperateService.ListByPeriod(start, end) // 获取平台收入的记录的map
	if err != nil {
		common.Logger.Warnln("ListByPeriod err ", err)
		common.Render(ctx, "27100203", err)
		return
	}

	alipayType, wechatType := 1, 2
	alipayMap, err := billService.ReportMapByPeriodAndAccountType(startAt[:10], endAt[:10], alipayType)
	if err != nil {
		common.Logger.Warnln("BasicByAccountTypeAndDatetime ------------", err)
		common.Render(ctx, "27100203", err)
		return
	}

	wechatMap, err := billService.ReportMapByPeriodAndAccountType(startAt[:10], endAt[:10], wechatType)
	if err != nil {
		common.Logger.Warnln("BasicByAccountTypeAndDatetime ------------", err)
		common.Render(ctx, "27100203", err)
		return
	}
	fileName := ""
	if startAt != "" && endAt != "" {
		if startAt[:10] != endAt[:10] {
			fileName = strings.Replace(startAt[5:10], "-", ".", -1) + "-" +
				strings.Replace(endAt[5:10], "-", ".", -1) + "结算报表"
		} else {
			fileName = strings.Replace(endAt[5:10], "-", ".", -1) + "结算报表"
		}
	} else {
		fileName = "结算报表"
	}
	tableHead := []interface{}{"日期", "平台收入-支付宝", "平台收入-微信", "批量付款金额-支付宝", "批量付款金额-微信", "手续费收入-支付宝", "手续费收入-微信", "手续费收入-总收入"}
	tableName := "结算报表"
	sheet, file, fileUrl, fileName, err := excel.GetExcelHeader(fileName, tableHead, tableName)
	if err != nil {
		common.Logger.Warningln("操作excel文件失败, err ------------>", err)
		common.Render(ctx, "27100204", err)
		return
	}

	for _, value := range *dailyOperateList {
		if excel.ExportBillReportDataAsCol(sheet, &value, (*alipayMap)[value.Date], (*wechatMap)[value.Date]) == 0 {
			common.Logger.Warningln("excel文件插入记录失败,err ------------>", err)
			common.Render(ctx, "27080505", err)
			return
		}
	}
	if err := file.Save(fileUrl); err != nil {
		common.Logger.Warningln("excel文件保存失败,err ------------>", err)
		common.Render(ctx, "27080506", err)
		return
	}
	sendFile := viper.GetString("server.href") + viper.GetString("export.loadsPath") + "/" + fileName
	common.Render(ctx, "27080500", map[string]string{"url": sendFile})
	return
}
