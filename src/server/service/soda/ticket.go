package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
)

type TicketService struct {
}

func (self *TicketService)TotalByDailyBillId(id int)(int,error){
	type Result struct {
		Total int
	}
	result := &Result{}
	r := common.SodaDB_R.Table("ticket").Where(" id = ?", "id").Count(result)
	if r.Error != nil {
		return nil,r.Error
	}
	return result.Total,nil
}

func (self *TicketService)DetailsByDailyBillId(id,limit,offset int)([]*soda.Ticket,error){
	ticket := []*soda.Ticket{}
	r := common.SodaDB_R.Where(" id = ? ",id).Limit(limit).Offset(offset).Find(ticket)
	if r.Error != nil {
		return nil,r.Error
	}
	return ticket,nil
}
