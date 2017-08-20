package service

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model"
)

type DailyBillService struct {
}

func (self *DailyBillService) TotalByBillId(billId string) (int, error) {
	var total int
	r := common.SodaMngDB_R.Model(&model.DailyBill{}).Where(" bill_id = ? ", billId).Count(&total)
	if r.Error != nil {
		return -1, r.Error
	}
	return total, nil
}

func (self *DailyBillService) ListByBillId(limit, offset int, billId string) ([]*model.DailyBill, error) {
	list := []*model.DailyBill{}
	params := make([]interface{}, 0)
	sql := "select * from daily_bill where daily_bill.deleted_at IS NULL "
	if billId != "" {
		sql += " and daily_bill.bill_id = ? "
		params = append(params, billId)
	}
	sql += " limit ? offset ? "
	params = append(params, limit)
	params = append(params, offset)
	common.Logger.Debugln("ListByBillId params===========", params)
	r := common.SodaMngDB_R.Raw(sql, params...).Order("id desc").Scan(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DailyBillService) BasicMapByBillAtAndUserId(billAt string, userIds interface{}) (map[int]*model.DailyBill, error) {
	list := &[]*model.DailyBill{}
	billMap := make(map[int]*model.DailyBill)
	r := common.SodaMngDB_R.Where("bill_at = ? and user_id in (?)", billAt, userIds).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, bill := range *list {
		billMap[bill.UserId] = bill
	}
	return billMap, nil
}

func (self *DailyBillService) BasicMap(billAt string, status int, userIds ...string) (map[int]*model.DailyBill, error) {
	list := &[]*model.DailyBill{}
	dailyBillMap := make(map[int]*model.DailyBill)
	r := common.SodaMngDB_R.Where("user_id in (?) and status = ? and bill_at = ?", userIds, status, billAt).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, dailyBill := range *list {
		dailyBillMap[dailyBill.UserId] = dailyBill
	}
	return dailyBillMap, nil
}

func (self *DailyBillService) BasicById(id int) (*model.DailyBill, error) {
	dailyBill := &model.DailyBill{}
	r := common.SodaMngDB_R.Where(" id = ? ", id).Find(dailyBill)
	if r.Error != nil {
		return nil, r.Error
	}
	return dailyBill, nil
}
