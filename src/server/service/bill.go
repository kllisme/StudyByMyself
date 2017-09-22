package service

import (
	"time"

	"github.com/go-errors/errors"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model"
)

type BillService struct {
}

// 有些函数需要连表查询
func (self *BillService) userIdByUserCondition(sql string, value []interface{}) ([]int, error) {
	ids := make([]int, 0)
	var id int
	r, err := common.SodaMngDB_R.Model(model.User{}).Select("id").Where(sql, value...).Rows()
	if err != nil {
		return ids, err
	}
	defer r.Close()
	for r.Next() {
		err = r.Scan(&id)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (self *BillService) TotalByAccountTypeAndTimeType(accountType, status, dateType int, startAt, endAt, keys string) (int, error) {
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
	if dateType == 1 {
		if startAt != "" {
			sql += " and Date(bill.created_at) >= ? "
			params = append(params, startAt[:10])
		}
		if endAt != "" {
			sql += " and Date(bill.created_at) <= ? "
			params = append(params, endAt[:10])
		}
	} else if dateType == 2 {
		if startAt != "" {
			sql += " and Date(bill.settled_at) >= ? "
			params = append(params, startAt[:10])
		}
		if endAt != "" {
			sql += " and Date(bill.settled_at) <= ? "
			params = append(params, endAt[:10])
		}
	}

	if keys != "" {
		userSql := " user.name like ? or user.account like ? "
		userId, err := self.userIdByUserCondition(userSql, []interface{}{"%" + keys + "%", "%" + keys + "%"})
		if err != nil {
			return -1, err
		}
		sql += " and bill.user_id in (?)"
		params = append(params, userId)
	}
	sql += " and bill.account_type = ? and user_id != 1 " // user_id != 1 过滤测试的账单
	params = append(params, accountType)
	r := common.SodaMngDB_R.Raw(sql, params...).Scan(result)
	if r.Error != nil {
		return -1, r.Error
	}
	return result.Total, nil

}
func (self *BillService) ListByAccountTypeAndTimeType(accountType, status, dateType, offset, limit int, startAt, endAt, keys string) ([]*model.Bill, error) {
	billList := []*model.Bill{}
	sql := "select * from bill where bill.deleted_at IS NULL "
	params := []interface{}{}
	if status > 0 {
		sql += " and bill.status = ? "
		params = append(params, status)
	}
	if dateType == 1 {
		if startAt != "" {
			sql += " and Date(bill.created_at) >= ? "
			params = append(params, startAt[:10])
		}
		if endAt != "" {
			sql += " and Date(bill.created_at) <= ? "
			params = append(params, endAt[:10])
		}
	} else if dateType == 2 {
		if startAt != "" {
			sql += " and Date(bill.settled_at) >= ? "
			params = append(params, startAt[:10])
		}
		if endAt != "" {
			sql += " and Date(bill.settled_at) <= ? "
			params = append(params, endAt[:10])
		}
	}
	// 运营商名称/登录账号
	if keys != "" {
		userSql := " user.name like ? or user.account like ? "
		userId, err := self.userIdByUserCondition(userSql, []interface{}{"%" + keys + "%", "%" + keys + "%"})
		if err != nil {
			return nil, err
		}
		sql += " and bill.user_id in (?) "
		params = append(params, userId)
	}
	sql += " and bill.account_type = ? and user_id != 1 " // user_id != 1 过滤测试的账单
	params = append(params, accountType)
	// 排序规则：等待结算的单排在最前面，然后到结算中→结算成功→结算失败，
	// 	同状态的账单按申请时间先后排序，最新提交的提现申请排在最前面；
	sql += " order by case " +
		"when bill.status=1 then 1 " +
		"when bill.status=3 then 2 " +
		"when bill.status=2 then 3 " +
		"when bill.status=4 then 4 " +
		"else 5 end asc, bill.id DESC "
	r := common.SodaMngDB_R.Raw(sql, params...).Limit(limit).Offset(offset).Scan(&billList)
	if r.Error != nil {
		return nil, r.Error
	}
	return billList, nil
}

func (self *BillService) ListByBillIdsAndStatus(billIds []interface{}, status []interface{}) ([]*model.Bill, error) {
	billList := []*model.Bill{}
	r := common.SodaMngDB_R.Where(" bill_id in (?) ", billIds).Where(" status in (?) ", status).Find(&billList)
	if r.Error != nil {
		return nil, r.Error
	}
	return billList, nil
}

func (self *BillService) BillTypeByBatchBill(billIds []interface{}) (int, error) {
	type Result struct {
		AccountType int
	}
	result := &Result{}
	accountType := -1
	sql := "select account_type from bill where bill.deleted_at IS NULL and bill_id = ? "
	for _, billId := range billIds {
		r := common.SodaMngDB_R.Raw(sql, billId).Scan(result)
		if r.Error != nil {
			return -1, r.Error
		}
		if accountType != -1 {
			if accountType != result.AccountType {
				return -1, errors.New("选取的账单存在不同的结算方式")
			}
		} else {
			accountType = result.AccountType
		}
	}
	return accountType, nil
}

func (self *BillService) BatchUpdateStatusById(status int, ids []interface{}) error {
	tx := common.SodaMngDB_WR.Begin()
	param := make(map[string]interface{}, 0)
	param["status"] = status
	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	// 接着更新daily_bill
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	tx.Commit()
	return nil
}

func (self *BillService) BatchUpdateStatusAndSettleAtById(status int, ids []interface{}) error {
	tx := common.SodaMngDB_WR.Begin()
	param := make(map[string]interface{}, 0)
	param["status"] = status // 只能2或者4
	// 先更新bill
	param["settled_at"] = time.Now()
	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	// 接着更新daily_bill
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	tx.Commit()
	return nil
}

func (self *BillService) BatchUpdateSubmitAtById(status int, ids []interface{}) error {
	tx := common.SodaMngDB_WR.Begin()
	param := make(map[string]interface{}, 0)
	param["status"] = status // 只能是3
	// 先更新bill
	//修改状态为"结账中",需更新"结账中时间"
	param["submitted_at"] = time.Now()
	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	dailyBillParam := make(map[string]interface{}, 0)
	// 接着更新daily_bill
	dailyBillParam["status"] = status
	dailyBillParam["submit_at"] = time.Now()
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ", ids).Update(dailyBillParam)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	tx.Commit()
	return nil
}

/*func (self *BillService) BatchUpdateStatusById(status int, ids []interface{}) error {
	tx := common.SodaMngDB_WR.Begin()
	param := make(map[string]interface{}, 0)
	param["status"] = status
	// 先更新bill
	if status == 2 || status == 4 {
		param["settled_at"] = time.Now()
	}
	if status == 3 {
		//修改状态为"结账中",需更新"结账中时间"
		param["submitted_at"] = time.Now()
	}
	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) ", ids).Update(param)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	dailyBillParam := make(map[string]interface{}, 0)
	// 接着更新daily_bill
	timeNow := time.Now().Local().Format("2006-01-02 15:04:05")
	if status == 2 || status == 4 {
		dailyBillParam["settled_at"] = timeNow
	}
	if status == 3 {
		dailyBillParam["submit_at"] = time.Now()
	}
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) ", ids).Update(dailyBillParam)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	tx.Commit()
	return nil
}*/

func (self *BillService) GetFirstWechatBill() (*model.Bill, error) {
	bill := &model.Bill{}
	// 找到状态为3结账中和结算账号为微信的账单出来结算
	r := common.SodaMngDB_R.Where(map[string]interface{}{"status": 3, "account_type": 2}).First(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}

func (self *BillService) Updates(list *[]*model.Bill) error {
	tx := common.SodaMngDB_WR.Begin()

	var successBillIDs []string
	var failBillIDs []string

	for _, _bill := range *list {
		if _bill.Status == 2 {
			successBillIDs = append(successBillIDs, _bill.BillId)
		} else if _bill.Status == 4 {
			failBillIDs = append(failBillIDs, _bill.BillId)
		}
	}

	r := tx.Model(&model.Bill{}).Where(" bill_id in (?) and status != 2 ", successBillIDs).Updates(map[string]interface{}{
		"status":     2,
		"settled_at": time.Now(),
	})
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	r = tx.Model(&model.Bill{}).Where(" bill_id in (?) and status != 2 ", failBillIDs).Updates(map[string]interface{}{
		"status":     4,
		"settled_at": time.Now(),
	})
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}

	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) and status != 2 ", successBillIDs).Updates(map[string]interface{}{
		"status":     2,
		"settled_at": time.Now(),
	})
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	r = tx.Model(&model.DailyBill{}).Where(" bill_id in (?) and status != 2 ", failBillIDs).Updates(map[string]interface{}{
		"status":     4,
		"settled_at": time.Now(),
	})
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}

	tx.Commit()
	return nil
}

func (self *BillService) BasicById(id int) (*model.Bill, error) {
	bill := &model.Bill{}
	r := common.SodaMngDB_R.Where(" id = ? ", id).Find(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}

func (self *BillService) BasicByBillId(billId string) (*model.Bill, error) {
	bill := &model.Bill{}
	r := common.SodaMngDB_R.Where(" bill_id = ? ", billId).Find(bill)
	if r.Error != nil {
		return nil, r.Error
	}
	return bill, r.Error
}

/* 返回日期字符串为key的map集合*/
func (self *BillService) ReportMapByPeriodAndAccountType(start, end string, accountType int) (*map[string]map[string]interface{}, error) {
	type Result struct {
		Cast int
		TotalAmount int
		SettledAt time.Time
	}

	sql := "select sum(cast) Cast,sum(total_amount) TotalAmount,date(settled_at) SettledAt from bill where " +
		"date(settled_at) >= ? and date(settled_at) <= ? and account_type = ? and status = 4 "+ // 必须要获取到成功的账单
		"group by date(SettledAt)"
	rows,err := common.SodaMngDB_R.Raw(sql,start,end,accountType).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	reportMap := make(map[string]map[string]interface{})
	for rows.Next() {
		result := &Result{Cast:0,TotalAmount:0}
		if err = rows.Scan(&result.Cast,&result.TotalAmount,&result.SettledAt);err != nil {
			return nil,err
		}
		resultMap := make(map[string]interface{})
		resultMap["cast"] = result.Cast
		resultMap["totalAmount"] = result.TotalAmount
		reportMap[result.SettledAt.Local().Format("2006-01-02")] = resultMap
	}
	return &reportMap, nil
}
