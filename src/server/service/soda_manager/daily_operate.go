package soda_manager

import (
	"time"

	"maizuo.com/soda/erp/api/src/server/common"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
)

type DailyOperateService struct {
}

/* 根据日期查找 */
func (self *DailyOperateService) BasicByDatetime(date string) (*mngModel.DailyOperate, error) { // TODO

	return nil, nil
}

func (self *DailyOperateService) ListByPeriod(startAt, endAt time.Time) (*[]mngModel.DailyOperate, error) {
	dailyOperateList := &[]mngModel.DailyOperate{}
	if r := common.SodaMngDB_R.Table("daily_operate").Where(" created_timestamp >= ?", startAt.Unix()).
		Where("created_timestamp < ? ", endAt.Unix()).Order("id,created_timestamp desc").Scan(dailyOperateList); r.Error != nil {
		return nil, r.Error
	}
	return dailyOperateList, nil
}

func (self *DailyOperateService) MapByPeriod(startAt, endAt time.Time) (*map[string]mngModel.DailyOperate, error) {
	dailyOperateList := &[]mngModel.DailyOperate{}
	if r := common.SodaMngDB_R.Table("daily_operate").Where(" created_timestamp >= ?", startAt.Unix()).
		Where("created_timestamp < ? ", endAt.Unix()).Order("id,created_timestamp desc").Scan(dailyOperateList); r.Error != nil {
		return nil, r.Error
	} else {
		m := make(map[string]mngModel.DailyOperate)
		for _, value := range *dailyOperateList {
			common.Logger.Debugln("value.Date ============== ", value.Date)
			m[value.Date] = value
		}
		return &m, nil
	}

}
