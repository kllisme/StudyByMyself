package soda

import (
	"maizuo.com/soda/erp/api/src/server/model/soda"
	"maizuo.com/soda/erp/api/src/server/common"
	"time"
	"maizuo.com/soda/erp/api/src/server/model"
	"github.com/jinzhu/gorm"
	"maizuo.com/soda/erp/api/src/server/entity"
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

func (self *TicketService) DetailsByDailyBill(dailyBill *model.DailyBill, limit, offset int) ([]*soda.Ticket, error) {
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

func (self *TicketService) Paging(userIDs []int, mobile string, paymentID int, deviceSerial string, deviceMode int, ownerIDs []int, feature int, appID int, statuses []int, start string, end string, offset int, limit int) (*entity.PaginationData, error) {
	pagination := entity.PaginationData{}
	ticketList := make([]*soda.Ticket, 0)
	db := common.SodaDB_R
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if len(userIDs) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("user_id in (?)", userIDs)
		})
	}

	if mobile != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("mobile = ?", mobile)
		})
	}

	if paymentID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("payment_id = ?", paymentID)
		})
	}

	if deviceSerial != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("device_serial like (?)", "%" + deviceSerial + "%")
		})
	}

	if deviceMode != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("device_mode = ?", deviceMode)
		})
	}

	if len(ownerIDs) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("owner_id in (?)", ownerIDs)
		})
	}

	if feature != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("feature = ?", feature)
		})
	}

	if appID != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("app_id = ?", appID)
		})
	}

	if start != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_timestamp >= UNIX_TIMESTAMP(?)", start)
		})
	}
	if end != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("created_timestamp < UNIX_TIMESTAMP(DATE_ADD(?, INTERVAL 1 day))", end)
		})
	}

	if len(statuses) != 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("status in (?)", statuses)
		})
	}

	if err := db.Model(&soda.Ticket{}).Scopes(scopes...).Count(&pagination.Pagination.Total).Offset(offset).Limit(limit).Order("id desc").Find(&ticketList).Error; err != nil {
		return nil, err
	}
	pagination.Pagination.From = offset + 1
	if limit == 0 {
		pagination.Pagination.To = pagination.Pagination.Total
	} else {
		pagination.Pagination.To = limit + offset
	}
	pagination.Objects = ticketList
	return &pagination, nil
}

func (self *TicketService) Refund(ticketID string) (*soda.Ticket, error) {
	tx := common.SodaDB_WR.Begin()
	_ticket := soda.Ticket{}
	if err := tx.Model(&soda.Ticket{}).Where("ticket_id = ? ", ticketID).Update("status", 4).Scan(&_ticket).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	_wallet := soda.Wallet{}
	if err := tx.Model(&soda.Wallet{}).Where("mobile = ? ", _ticket.Mobile).First(&_wallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&soda.Wallet{}).Where(_wallet.ID).Update("value", _wallet.Value + _ticket.Value).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&soda.Bill{
		Mobile:_wallet.Mobile,
		UserID:_wallet.UserID,
		WalletID:_wallet.ID,
		Title:"退款",
		Value:_ticket.Value,
		Type:4,
		Action:1,
		Status:0,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &_ticket, nil
}

func (self *TicketService) GetByTicketID(ticketID string) (*soda.Ticket, error) {
	ticket := soda.Ticket{}
	err := common.SodaDB_R.Where("ticket_id = ?", ticketID).Find(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
