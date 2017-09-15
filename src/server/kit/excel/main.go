package excel

import (
	"github.com/tealeg/xlsx"
	"os"
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/payload/crm"
)
// 生成excel表头,并返回文件路径以及名字
func GetExcelHeader(fileName string, values []interface{},tableName string)(sheet *xlsx.Sheet,file *xlsx.File, url string, name string, err error){
	root, _ := os.Getwd()
	path := root + "/temp"
	name = fileName + ".xlsx"
	url = path + "/" + name
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(tableName)
	if err != nil {
		return nil,nil, "", "", err
	}
	sheet.AddRow().WriteSlice(&values, -1)
	return
}
// 添加一行数据
func ExportBillDataAsCol(sheet *xlsx.Sheet, bill *model.Bill) (int) {
	row := sheet.AddRow()
	mode := "自动结算"
	if bill.Mode != 0 { // 0代表自动提现
		mode = "手动结算"
	}
	status := "等待结算"
	settledAt := "-"
	switch bill.Status {
	case 1:
		status = "等待结算"
	case 2:
		status = "结算成功"
		settledAt = bill.SettledAt.Local().Format("2006-01-02 15:04")
	case 3:
		status = "结算中"
	case 4:
		status = "结算失败"
		settledAt = bill.SettledAt.Local().Format("2006-01-02 15:04")
	}
	s := []interface{}{
		bill.CreatedAt.Local().Format("2006-01-02 15:04"),
		bill.UserName+"|"+bill.UserAccount,
		bill.RealName+"|账号:"+bill.Account,
		bill.BillId,
		bill.Count,
		float64(bill.TotalAmount)/100.00,
		float64(bill.Cast)/100.00,
		float64(bill.Amount)/100.00,
		status,
		settledAt,
		mode,
	}
	return row.WriteSlice(&s, -1)
}

func ExportConsumptionAsCol(sheet *xlsx.Sheet, consumption *crm.Consumption) (int) {
	row := sheet.AddRow()
	s := []interface{}{
		consumption.TicketID,
		consumption.ParentOperator +" ("+consumption.ParentOperatorMobile + ")",
		consumption.Operator,
		consumption.Telephone,
		consumption.DeviceSerial,
		consumption.Address,
		consumption.CustomerMobile,
		consumption.Password,
		consumption.TypeName,
		float64(consumption.Value)/100.00,
		consumption.Payment,
		consumption.CreatedAt.Local().Format("2006-01-02 15:04"),
	}
	return row.WriteSlice(&s, -1)
}
