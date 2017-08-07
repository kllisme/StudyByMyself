package service

import (
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/common"
	"time"
)

type BillService struct {

}

func (self *BillService) TotalByAccountType(accountType,status int,createdAt,settledAt,keys string)(int,error){
	type Result struct {
		Total int
	}
	result := &Result{}
	sql := "select count(*) as total from bill where bill.deleted_at IS NULL "
	params := []interface{}{}
	if status > 0 {
		sql += " and bill.status = ? "
		params = append(params, status)
	}

	if createdAt != "" {
		sql += " and Date(bill.created_at) = ? "
		params = append(params, createdAt)
	}
	if settledAt != "" {
		sql += " and Date(bill.settled_at) = ? "
		params = append(params, settledAt)
	}

	if keys != "" {
		sql += " and (bill.user_name like ? or bill.account_name like ?) "
		params = append(params, "%"+keys+"%")
		params = append(params, "%"+keys+"%")
	}
	sql += " and bill.account_type = ? and user_id != 1 " // user_id != 1 过滤测试的账单
	params = append(params,accountType)
	r := common.SodaMngDB_R.Raw(sql, params...).Scan(result)
	if r.Error != nil {
		return -1,r.Error
	}
	return result.Total,nil
}

func (self *BillService) ListByAccountType(accountType,status,offset,limit int,createdAt,settledAt,keys string )([]*model.Bill,error){
	billList := []*model.Bill{}
	sql := "select * from bill where ( bill.deleted_at IS NULL "
	params := []interface{}{}
	if status > 0 {
		sql += " and bill.status = ? "
		params = append(params, status)
	}
	if createdAt != "" {
		sql += " and Date(bill.created_at) = ? "
		params = append(params, createdAt)
	}
	if settledAt != "" {
		sql += " and Date(bill.settled_at) = ? "
		params = append(params, settledAt)
	}
	params = append(params, createdAt)
	if keys != "" {
		sql += " and (bill.user_name like ? or bill.account_name like ?) "
		params = append(params, "%"+keys+"%")
		params = append(params, "%"+keys+"%")
	}
	sql += " and bill.account_type = ? and user_id != 1 " // user_id != 1 过滤测试的账单
	params = append(params,accountType)
	// TODO 排序规则：等待结算的单排在最前面，然后到结算中→结算成功→结算失败，
	// 	同状态的账单按申请时间先后排序，最新提交的提现申请排在最前面；
	r := common.SodaMngDB_R.Raw(sql, params...).Order(" created_at desc ").Offset(offset).Limit(limit).Scan(&billList)
	if r.Error != nil {
		return nil,r.Error
	}
	return billList,nil
}

func (self *BillService) ListByBillIdsAndStatus(billIds []string,status []int)([]*model.Bill,error){
	billList := []*model.Bill{}
	r := common.SodaMngDB_R.Where(" bill_id in (?) ",billIds...).Where(" status in (?) ",status...).Find(billList)
	if r.Error != nil {
		return nil,r.Error
	}
	return billList,nil
}

func (self *BillService) BatchUpdateStatusById(status int, ids ...interface{}) (error) {
	tx := common.SodaMngDB_R.Begin()
	param := make(map[string]interface{}, 0)
	param["status"] = status
	// 先更新bill
	if status == 2 {
		param["settled_at"] = time.Now()
	}
	if status == 3 {
		//修改状态为"结账中",需更新"结账中时间"
		param["submitted_at"] = time.Now()
	}
	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) ", ids...).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	dailyBillParam := make(map[string]interface{}, 0)
	// 接着更新daily_bill
	timeNow := time.Now().Local().Format("2006-01-02 15:04:05")
	if status == 2 {
		dailyBillParam["settled_at"] = timeNow
	}
	if status == 3 {
		dailyBillParam["submit_at"] = time.Now()
	}
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ",ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	return  nil
}
