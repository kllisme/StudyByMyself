package soda_manager

import (
	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type DailyBillService struct {
}

func (self *DailyBillService) TotalByBillId(billId string) (int, error) {
	var total int
	r := common.SodaMngDB_R.Model(&mngModel.DailyBill{}).Where(" bill_id = ? ", billId).Count(&total)
	if r.Error != nil {
		return -1, r.Error
	}
	return total, nil
}

func (self *DailyBillService) ListByBillId(limit, offset int, billId string) ([]*mngModel.DailyBill, error) {
	list := []*mngModel.DailyBill{}
	params := make([]interface{}, 0)
	sql := "select * from daily_bill where daily_bill.deleted_at IS NULL "
	if billId != "" {
		sql += " and daily_bill.bill_id = ? "
		params = append(params, billId)
	}
	sql += " order by bill_at desc,id desc limit ? offset ?"
	params = append(params, limit)
	params = append(params, offset)
	common.Logger.Debugln("ListByBillId params===========", params)
	r := common.SodaMngDB_R.Raw(sql, params...).Scan(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DailyBillService) BasicMapByBillAtAndUserId(billAt string, userIds interface{}) (map[int]*mngModel.DailyBill, error) {
	list := &[]*mngModel.DailyBill{}
	billMap := make(map[int]*mngModel.DailyBill)
	r := common.SodaMngDB_R.Where("bill_at = ? and user_id in (?)", billAt, userIds).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, bill := range *list {
		billMap[bill.UserId] = bill
	}
	return billMap, nil
}

func (self *DailyBillService) BasicMap(billAt string, status int, userIds ...string) (map[int]*mngModel.DailyBill, error) {
	list := &[]*mngModel.DailyBill{}
	dailyBillMap := make(map[int]*mngModel.DailyBill)
	r := common.SodaMngDB_R.Where("user_id in (?) and status = ? and bill_at = ?", userIds, status, billAt).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, dailyBill := range *list {
		dailyBillMap[dailyBill.UserId] = dailyBill
	}
	return dailyBillMap, nil
}

func (self *DailyBillService) BasicById(id int) (*mngModel.DailyBill, error) {
	dailyBill := &mngModel.DailyBill{}
	r := common.SodaMngDB_R.Where(" id = ? ", id).Find(dailyBill)
	if r.Error != nil {
		return nil, r.Error
	}
	return dailyBill, nil
}
