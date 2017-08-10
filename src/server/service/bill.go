package service

import (
	"time"

	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model"
	"github.com/go-errors/errors"
)

type BillService struct {
}

func (self *BillService) TotalByAccountType(accountType, status int, createdAt, settledAt, keys string) (int, error) {
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
	params = append(params, accountType)
	r := common.SodaMngDB_R.Raw(sql, params...).Scan(result)
	if r.Error != nil {
		return -1, r.Error
	}
	return result.Total, nil
}

func (self *BillService) ListByAccountType(accountType, status, offset, limit int, createdAt, settledAt, keys string) ([]*model.Bill, error) {
	billList := []*model.Bill{}
	sql := "select * from bill where bill.deleted_at IS NULL "
	params := []interface{}{}
	if status > 0 {
		sql += " and bill.status = ? "
		params = append(params, status)
	}
	if createdAt != "" {
		sql += " and Date(bill.created_at) = ? "
		params = append(params, []byte(createdAt)[0:10])
	}
	if settledAt != "" {
		sql += " and Date(bill.settled_at) = ? "
		params = append(params, []byte(settledAt)[0:10])
	}
	if keys != "" {
		sql += " and (bill.user_name like ? or bill.account_name like ?) "
		params = append(params, "%"+keys+"%")
		params = append(params, "%"+keys+"%")
	}
	sql += " and bill.account_type = ? and user_id != 1 " // user_id != 1 过滤测试的账单
	params = append(params, accountType)
	common.Logger.Debugln("params : ",params)
	// 排序规则：等待结算的单排在最前面，然后到结算中→结算成功→结算失败，
	// 	同状态的账单按申请时间先后排序，最新提交的提现申请排在最前面；
	sql += " order by case " +
		"when bill.status=1 then 1 " +
		"when bill.status=3 then 2 " +
		"when bill.status=2 then 3 " +
		"when bill.status=4 then 4 " +
		"else 5 end asc, bill.created_at DESC "
	r := common.SodaMngDB_R.Raw(sql, params...).Limit(limit).Offset(offset).Order(" created_at desc ").Scan(&billList)
	if r.Error != nil {
		return nil, r.Error
	}
	return billList, nil
}

func (self *BillService) ListByBillIdsAndStatus(billIds []interface{}, status []interface{}) ([]*model.Bill, error) {
	billList := []*model.Bill{}
	r := common.SodaMngDB_R.Where(" bill_id in (?) ", billIds...).Where(" status in (?) ", status...).Find(billList)
	if r.Error != nil {
		return nil, r.Error
	}
	return billList, nil
}

func(self *BillService) BillTypeByBatchBill(billIds []interface{}) (int, error) {
	type Result struct {
		AccountType int
	}
	result := &Result{}
	accountType := -1
	sql := "select account_type as `AccountType` from bill where bill.deleted_at IS NULL and bill_id = ? "
	for _,billId := range billIds {
		r := common.SodaMngDB_R.Where(sql, billId).Find(result)
		if r.Error != nil {
			return -1, r.Error
		}
		if accountType != 0 {
			if accountType != result.AccountType {
				return -1,errors.New("选取的账单存在不同的结算方式")
			}
		}else{
			accountType = result.AccountType
		}
	}
	return accountType, nil
}

func (self *BillService) BatchUpdateStatusById(status int, ids ...interface{}) error {
	tx := common.SodaMngDB_WR.Begin()
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
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	return nil
}

func (self *BillService) GetFirstWechatBill() (*model.Bill, error) {
	bill := &model.Bill{}
	// 找到状态为3结账中和结算账号为微信的账单出来结算
	r := common.SodaMngDB_R.Where(map[string]interface{}{"status": 3, "account_type": 2}).First(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}


func (self *BillService) Updates(list *[]*model.Bill) (int, error) {
	var rows int
	rows = 0
	tx := common.SodaMngDB_WR.Begin()
	for _, _bill := range *list {
		param := map[string]interface{}{
			"status":     _bill.Status,
			"settled_at": _bill.SettledAt,
		}
		r := tx.Model(&model.Bill{}).Where(" id = ? and status != 2 ", _bill.ID).Updates(param)
		if r.Error != nil {
			common.Logger.Warningln(r.Error.Error())
			tx.Rollback()
			return 0, r.Error
		} else {
			rows += int(r.RowsAffected)
		}
		r = tx.Model(&model.DailyBill{}).Where(" bill_id = ? and status != 2",_bill.BillId).Updates(param)
		if r.Error != nil {
			common.Logger.Warningln(r.Error.Error())
			tx.Rollback()
			return 0, r.Error
		} else {
			rows += int(r.RowsAffected)
		}
	}
	tx.Commit()
	return rows, nil
}

func (self *BillService)BasicById(id int)(*model.Bill,error){
	bill := &model.Bill{}
	r := common.SodaMngDB_R.Where(" id = ? ",id).Find(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}

func (self *BillService)BasicByBillId(billId string)(*model.Bill,error){
	bill := &model.Bill{}
	r := common.SodaMngDB_R.Where(" bill_id = ? ",billId).Find(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}
