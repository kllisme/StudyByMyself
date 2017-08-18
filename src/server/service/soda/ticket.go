package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
	"time"
	"maizuo.com/soda/erp/api/src/server/model"
)

type TicketService struct {
}


func (self *TicketService) TotalByDailyBill(dailyBill *model.DailyBill) (int, error) {
	t, _ := time.Parse("2006-01-02T15:04:05+08:00", dailyBill.BillAt)
	tomorrow := t.Local().AddDate(0, 0, 1).Format("2006-01-02")
	var total int64 = 0
	sql := `
	owner_id=convert(?,signed)
	and
	created_timestamp>=unix_timestamp(?)
	and
	created_timestamp<unix_timestamp(?)
	and
	status in (4,7)
	`
	r := common.SodaDB_R.Model(&soda.Ticket{}).Where(sql, dailyBill.UserId, dailyBill.BillAt, tomorrow).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return int(total), nil
}

func (self *TicketService)DetailsByDailyBill(dailyBill *model.DailyBill,limit,offset int)([]*soda.Ticket,error){
	list := []*soda.Ticket{}
	//t, _ := time.Parse("2006-01-02", dailyBill.BillAt)
	t, _ := time.Parse("2006-01-02T15:04:05+08:00", dailyBill.BillAt)
	tomorrow := t.Local().AddDate(0, 0, 1).Format("2006-01-02")
	sql := `
	owner_id=convert(?,signed)
	and
	created_timestamp>=unix_timestamp(?)
	and
	created_timestamp<unix_timestamp(?)
	and
	status in (4,7)
	`
	r := common.SodaDB_R.Model(&soda.Ticket{}).Where(sql, dailyBill.UserId, dailyBill.BillAt, tomorrow).Offset(offset).Limit(limit).Find(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}
