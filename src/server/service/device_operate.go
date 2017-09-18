package service

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	"maizuo.com/soda/erp/api/src/server/model"
	"strings"
)

type DeviceOperateService struct {
}

func (self *DeviceOperateService) Create(deviceOperate *model.DeviceOperate) (int, error) {
	var err error
	isTrue := common.SodaMngDB_WR.NewRecord(deviceOperate)
	if !isTrue {
		e := &functions.DefinedError{}
		e.Msg = "can not create a new record!"
		err = e
		return 0, err
	}
	r := common.SodaMngDB_WR.Create(&deviceOperate)
	if r.Error != nil {
		return 0, r.Error
	}
	return int(r.RowsAffected), nil
}

func (self *DeviceOperateService) BatchCreate(deviceOperates []*model.DeviceOperate) (int, error) {
	params := make([]string, 0)
	values := make([]interface{}, 0)
	if len(deviceOperates) <= 0 {
		return 0, error(&functions.DefinedError{"length of deviceOperares no bigger than 0"})
	}
	sql := "insert into `device_operate` (`operator_id`,`operator_type`,`serial_number`,`user_id`,`from_user_id`,`to_user_id`,`created_at`,`updated_at`) values"
	for _, operate := range deviceOperates {
		params = append(params, "(?,?,?,?,?,?,now(),now())")
		values = append(values, operate.OperatorID, operate.OperatorType, operate.SerialNumber, operate.UserID, operate.FromUserID, operate.ToUserID)
	}
	sql += strings.Join(params, ",")
	result := common.SodaMngDB_WR.Exec(sql, values...)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected <= 0 {
		return int(result.RowsAffected), error(&functions.DefinedError{"RowsAffected no bigger than 0"})
	}
	return int(result.RowsAffected), nil
}

func (self *DeviceOperateService) BatchCreateBussiness(operatorId int, operatorType int, deviceMap map[string]*model.Device, serialNums []string, toUserId int) (int, error) {
	deviceOperates := make([]*model.DeviceOperate, 0)
	for _, s := range serialNums {
		_device := &model.Device{}
		if deviceMap[s] == nil {
			continue
		}
		_device = deviceMap[s]
		deviceOperate := &model.DeviceOperate{
			OperatorID:   operatorId,
			OperatorType: operatorType,
			SerialNumber: s,
			UserID:       _device.UserID,
			FromUserID:   _device.FromUserID,
			ToUserID:     toUserId,
		}
		deviceOperates = append(deviceOperates, deviceOperate)
	}
	return self.BatchCreate(deviceOperates)
}
