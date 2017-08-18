package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model"
)

type DeviceService struct {
}

func (self *DeviceService) Basic(id int) (*model.Device, error) {
	device := &model.Device{}
	r := common.SodaMngDB_R.Where("id = ?", id).First(device)
	if r.Error != nil {
		return nil, r.Error
	}
	return device, nil
}

func (self *DeviceService) BasicBySerialNumber(serialNumber string) (*model.Device, error) {
	device := &model.Device{}
	r := common.SodaMngDB_R.Where("serial_number = ?", serialNumber).First(device)
	if r.Error != nil {
		return nil, r.Error
	}
	return device, nil
}

func (self *DeviceService) BasicMapBySerialNumber(serialNumber ...string) (map[string]*model.Device, error) {
	list := &[]*model.Device{}
	data := make(map[string]*model.Device, 0)
	r := common.SodaMngDB_R.Where("serial_number in (?)", serialNumber).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	for _, _device := range *list {
		data[_device.SerialNumber] = _device
	}
	return data, nil
}

func (self *DeviceService) ListBySerialNumber(userId int, schoolId int, serialNumber ...string) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	params := make([]interface{}, 0)
	sql := "user_id = ? and "
	params = append(params, userId)
	if schoolId != -1 {
		sql += " school_id = ? and "
		params = append(params, schoolId)
	}
	sql += " serial_number in (?)"
	params = append(params, serialNumber)
	r := common.SodaMngDB_R.Where(sql, params...).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) List(page int, perPage int) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	r := common.SodaMngDB_R.Offset((page - 1) * perPage).Limit(perPage).Order("id desc").Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) Total() (int, error) {
	var total int
	r := common.SodaMngDB_R.Model(&model.Device{}).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return total, nil
}

func (self *DeviceService) TotalByUser(userId int) (int, error) {
	device := &model.Device{}
	var total int
	r := common.SodaMngDB_R.Model(device).Where("user_id = ?", userId).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return total, nil
}

func (self *DeviceService) TotalByUserAndNextLevel(isAssigned bool, user *model.User, serialNumber string, userQuery string, isEqual bool) (int, error) {
	var total int
	params := make([]interface{}, 0)

	sql := ""
	if userQuery != "" {
		sql = " select count(*) as total from device d, user u where 1=1 "
	} else {
		sql = " select count(*) as total from device d where 1=1 "
	}
	/*if user.ID == 1 {
		sql += " and (user_id = ? or from_user_id= ? or has_retrofited = 1) "
	} else {
		sql += " and (user_id = ? or from_user_id= ? and has_retrofited = 0) "
	}
	params = append(params, user.ID, user.ID)*/
	if isAssigned {
		//分配出去的设备
		if user.ID == 1 {
			sql += " and (from_user_id = ? or has_retrofited = 1) "
		} else {
			sql += " and (from_user_id = ? and has_retrofited = 0) "
		}
		params = append(params, user.ID)
	} else {
		//自己的设备
		if user.ID == 1 {
			sql += " and (user_id = ? or has_retrofited = 1) "
		} else {
			sql += " and (user_id = ? and has_retrofited = 0) "
		}
		params = append(params, user.ID)
	}
	if serialNumber != "" {
		if isEqual {
			sql += " and d.serial_number in (" + serialNumber + ") "
		} else {
			sql += " and d.serial_number like ? "
			params = append(params, "%"+serialNumber+"%")
		}
	} else if userQuery != "" {
		sql += " and ( (u.ID=d.user_id /*or u.ID=d.from_user_Id*/) and (u.name like ? or u.account like ? or u.contact like ? ) ) "
		params = append(params, "%"+userQuery+"%", "%"+userQuery+"%", "%"+userQuery+"%")
	}
	r := common.SodaMngDB_R.Raw(sql, params...).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return total, nil
}

func (self *DeviceService) ListByUserAndNextLevel(isAssigned bool, user *model.User, serialNumber string, userQuery string, isEqual bool, page int, perPage int) (*[]*model.Device, error) {
	list := make([]*model.Device, 0)
	params := make([]interface{}, 0)
	sql := ""
	if userQuery != "" {
		sql = " select u.name as user_name, d.* from device d, user u where 1=1 "
	} else {
		sql = " select d.* from device d where 1=1 "
	}
	/*if user.ID == 1 {
		sql += " and (user_id = ? or from_user_id= ? or has_retrofited = 1) "
	} else {
		sql += " and (user_id = ? or from_user_id= ? and has_retrofited = 0) "
	}
	params = append(params, user.ID, user.ID)*/
	if isAssigned {
		//分配出去的设备
		if user.ID == 1 {
			sql += " and (from_user_id = ? or has_retrofited = 1) "
		} else {
			sql += " and (from_user_id = ? and has_retrofited = 0) "
		}
		params = append(params, user.ID)
	} else {
		//自己的设备
		if user.ID == 1 {
			sql += " and (user_id = ? or has_retrofited = 1) "
		} else {
			sql += " and (user_id = ? and has_retrofited = 0) "
		}
		params = append(params, user.ID)
	}
	if serialNumber != "" {
		if isEqual {
			sql += " and d.serial_number in (" + serialNumber + ") "
		} else {
			sql += " and d.serial_number like ? "
			params = append(params, "%"+serialNumber+"%")
		}
	} else if userQuery != "" {
		sql += " and ( (u.ID=d.user_id ) and (u.name like ? or u.account like ? or u.contact like ? ) ) "
		params = append(params, "%"+userQuery+"%", "%"+userQuery+"%", "%"+userQuery+"%")
	}
	//测试按设备编号排序
	if user.ID == 1 {
		sql += " order by case when user_id=? then 1 else 2 end asc, serial_number desc "
	} else {
		//其它用户按楼层排序
		sql += " order by case when user_id=? then 1 else 2 end asc, address desc, id desc "
	}
	params = append(params, user.ID)
	if perPage > 0 && page > 0 {
		sql += " limit ? offset ? "
		params = append(params, perPage, (page-1)*perPage)
	}
	rows, err := common.SodaMngDB_R.Raw(sql, params...).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		device := &model.Device{}
		common.SodaMngDB_R.ScanRows(rows, device)
		list = append(list, device)

	}
	return &list, nil
}

func (self *DeviceService) ListByUser(userId int, page int, perPage int) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	//limit perPage offset (page-1)*perPage
	r := common.SodaMngDB_R.Offset((page-1)*perPage).Limit(perPage).Where("user_id = ?", userId).Order("id desc").Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) ListByUserOrderByAddress(userId int, page int, perPage int) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	r := common.SodaMngDB_R.Offset((page-1)*perPage).Limit(perPage).Where("user_id = ?", userId).Order("address desc,serial_number desc").Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) ListBySerialNumbers(serialNumbers ...string) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	r := common.SodaMngDB_R.Where("serial_number in (?) ", serialNumbers).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) ListByUserAndSerialNumbers(userId int, serialNumbers []string) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	r := common.SodaMngDB_R.Where("user_id = ? and serial_number in (?) ", userId, serialNumbers).Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) ListByUserAndSchool(userId int, schoolId int, page int, perPage int, serialNums ...string) (*[]*model.Device, error) {
	list := &[]*model.Device{}
	params := make([]interface{}, 0)
	sql := "user_id = ? and school_id = ?"
	params = append(params, userId, schoolId)
	if len(serialNums) > 0 && serialNums[0] != "" {
		sql += " and serial_number in (?)"
		params = append(params, serialNums)
	}
	r := common.SodaMngDB_R.Offset((page-1)*perPage).Limit(perPage).Where(sql, params...).Order("id desc").Find(list)
	if r.Error != nil {
		return nil, r.Error
	}
	return list, nil
}

func (self *DeviceService) TotalByByUserAndSchool(userId int, schoolId int, serialNums ...string) (int, error) {
	device := &model.Device{}
	var total int
	params := make([]interface{}, 0)
	sql := "user_id = ? and school_id = ?"
	params = append(params, userId, schoolId)

	if len(serialNums) > 0 && serialNums[0] != "" {
		sql += " and serial_number in (?)"
		params = append(params, serialNums)
	}
	r := common.SodaMngDB_R.Model(device).Where(sql, params...).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return total, nil
}

func (self *DeviceService) Create(device *model.Device) bool {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Create(device).Scan(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Create BoxInfo:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) Update(device model.Device, isBatchUpdate bool) bool {
	var r *gorm.DB
	transAction := common.SodaMngDB_WR.Begin()
	var deviceMap = make(map[string]interface{})
	if isBatchUpdate {
		if device.FirstPulsePrice >= 0 && device.SecondPulsePrice >= 0 && device.ThirdPulsePrice >= 0 && device.FourthPulsePrice >= 0 {
			deviceMap["first_pulse_price"] = device.FirstPulsePrice
			deviceMap["second_pulse_price"] = device.SecondPulsePrice
			deviceMap["third_pulse_price"] = device.ThirdPulsePrice
			deviceMap["fourth_pulse_price"] = device.FourthPulsePrice
		}
		if device.FirstPulseName != "" && device.SecondPulseName != "" && device.ThirdPulseName != "" && device.FourthPulseName != "" {
			deviceMap["first_pulse_name"] = device.FirstPulseName
			deviceMap["second_pulse_name"] = device.SecondPulseName
			deviceMap["third_pulse_name"] = device.ThirdPulseName
			deviceMap["fourth_pulse_name"] = device.FourthPulseName
		}
		if device.ProvinceId >= 0 && device.SchoolId >= 0 {
			deviceMap["province_id"] = device.ProvinceId
			deviceMap["school_id"] = device.SchoolId
		}
		if device.Address != "" {
			deviceMap["address"] = device.Address
		}
	} else {
		deviceMap["first_pulse_price"] = device.FirstPulsePrice
		deviceMap["second_pulse_price"] = device.SecondPulsePrice
		deviceMap["third_pulse_price"] = device.ThirdPulsePrice
		deviceMap["fourth_pulse_price"] = device.FourthPulsePrice
		deviceMap["first_pulse_name"] = device.FirstPulseName
		deviceMap["second_pulse_name"] = device.SecondPulseName
		deviceMap["third_pulse_name"] = device.ThirdPulseName
		deviceMap["fourth_pulse_name"] = device.FourthPulseName
		deviceMap["city_id"] = device.CityId
		deviceMap["province_id"] = device.ProvinceId
		deviceMap["school_id"] = device.SchoolId
		deviceMap["address"] = device.Address
		deviceMap["label"] = device.Label
		deviceMap["reference_device_id"] = device.ReferenceDeviceId
	}
	serialNumber := []string{device.SerialNumber}
	if isBatchUpdate {
		serialNumber = strings.Split(device.SerialNumber, ",")
		r = transAction.Model(&model.Device{}).Where("serial_number in (?)", serialNumber).Updates(deviceMap).Scan(&device)
	} else {
		r = transAction.Model(&model.Device{}).Where("id = ?", device.ID).Updates(deviceMap).Scan(&device)
	}
	if r.Error != nil {
		common.Logger.Warningln("DB Update Device:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) UpdatePulseName(device model.Device) bool {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Model(&model.Device{}).Where("id = ?", device.ID).Updates(device)
	if r.Error != nil {
		common.Logger.Warningln("DB UpdatePulseName Device:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) UpdateStatus(device model.Device) bool {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Model(&model.Device{}).Where("id = ?", device.ID).Update("status", device.Status).Scan(&device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-STATUS:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) UpdateBySerialNumber(device *model.Device) bool {
	transAction := common.SodaMngDB_WR.Begin()
	r := transAction.Model(&model.Device{}).Where("serial_number = ?", device.SerialNumber).Updates(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-BY-DEVICENO:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	//再单独更新一次学校id因为传入的学校id可以为0
	r = transAction.Model(&model.Device{}).Where("serial_number = ?", device.SerialNumber).Update("school_id", device.SchoolId).Scan(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-BY-DEVICENO:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	//单独更新一次学校id
	transAction.Commit()
	return true
}

func (self *DeviceService) BatchUpdateBySerialNumber(device *model.Device, serialList []string) bool {
	transAction := common.SodaMngDB_WR.Begin()
	device.SerialNumber = "" //不更新这个字段
	r := transAction.Model(&model.Device{}).Where("serial_number IN (?)", serialList).Updates(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-BY-DEVICENO:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	//再单独更新一次可能为0的值，学校id，4个单洗价格
	var zeroValue = make(map[string]interface{})
	zeroValue["school_id"] = device.SchoolId
	zeroValue["first_pulse_price"] = device.FirstPulsePrice
	zeroValue["second_pulse_price"] = device.SecondPulsePrice
	zeroValue["third_pulse_price"] = device.ThirdPulsePrice
	zeroValue["fourth_pulse_price"] = device.FourthPulsePrice
	r = transAction.Model(&model.Device{}).Where("serial_number IN (?)", serialList).Updates(zeroValue).First(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-BY-DEVICENO:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) Reset(id int, user *model.User) bool {
	device := &model.Device{}
	transAction := common.SodaMngDB_WR.Begin()
	data := map[string]interface{}{
		"user_id":        1,
		"from_user_id":   user.ID,
		"status":         9,
		"has_retrofited": 1,
	}
	r := transAction.Model(&model.Device{}).Where("id = ?", id).Updates(data).Scan(device)
	if r.Error != nil {
		common.Logger.Warningln("DB Update BoxInfo-Reset:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) ResetPasswordStep(serialNumber string, step int) error {
	transAction := common.SodaMngDB_WR.Begin()
	data := map[string]interface{}{
		"step": step,
	}
	r := transAction.Model(&model.Device{}).Where("serial_number = ?", serialNumber).Updates(data)
	if r.Error != nil {
		transAction.Rollback()
		return r.Error
	}
	transAction.Commit()
	return nil
}

func (self *DeviceService) Delete(id int) bool {
	device := &model.Device{}
	transAction := common.SodaMngDB_WR.Begin()
	//硬删除记录
	r := transAction.Model(&model.Device{}).Where("id = ?", id).Scan(device).Unscoped().Delete(&model.Device{})
	if r.Error != nil {
		common.Logger.Warningln("DB Delete BoxInfo:", r.Error.Error())
		transAction.Rollback()
		return false
	}
	transAction.Commit()
	return true
}

func (self *DeviceService) ListSchoolIdByUser(userId int) (*[]int, error) {
	var schoolList []int
	type MyDevice struct {
		SchoolId int `json:"school_id"`
	}
	lists := &[]*MyDevice{}
	//带去重
	r := common.SodaMngDB_R.Raw("SELECT DISTINCT school_id FROM device WHERE user_id = ? AND deleted_at IS NULL", userId).Scan(lists)
	if r.Error != nil {
		return nil, r.Error
	}
	//赋值到数组
	for _, list := range *lists {
		schoolList = append(schoolList, list.SchoolId)
	}
	return &schoolList, nil
}

//通过用户列表计算一个有多少个设备
func (self *DeviceService) TotalByUserIds(userIds []int) (int, error) {
	//list := &[]*model.Device{}
	var total int64
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id IN (?)", userIds).Count(&total)
	if r.Error != nil {
		return 0, r.Error
	}
	return int(total), nil
}

//判断这批设备序列号是否都是属于test用户或者自身用户的
func (self *DeviceService) IsOwntoMeOrTest(userId int, serialList []string) bool {
	var total int64
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id IN (1,?) AND serial_number IN (?)", userId, serialList).Count(&total)
	if r.Error != nil {
		common.Logger.Warningln(r.Error)
		return false
	}
	if int(total) == len(serialList) {
		return true
	} else {
		return false
	}
}

//判断这一批序列号要不还没注册要不还没录入
func (self *DeviceService) IsNoExist(serialList []string) bool {
	var total int64
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("serial_number IN (?)", serialList).Count(&total)
	if r.Error != nil {
		common.Logger.Warningln(r.Error)
		return false
	}
	if total == 0 {
		return true
	} else {
		return false
	}
}

func (self *DeviceService) BatchCreateBySerialNum(device *model.Device, serialList []string) error {
	transAction := common.SodaMngDB_WR.Begin()
	device.UserId = 1 //添加的设备userid为1
	device.CreatedAt = time.Now()
	step := viper.GetInt("device.step")
	device.Step = step
	sql := "INSERT INTO device (user_id,label,step,serial_number,reference_device_id,province_id,city_id,district_id,school_id,address,first_pulse_price,second_pulse_price,third_pulse_price,fourth_pulse_price,first_pulse_name,second_pulse_name,third_pulse_name,fourth_pulse_name,status,created_at,updated_at) VALUES "
	for k, serial := range serialList {
		val := fmt.Sprintf("(%d,'%s',%d,'%s',%d,%d,%d,%d,%d,'%s',%d,%d,%d,%d,'%s','%s','%s','%s',%d,'%s','%s')", device.UserId, device.Label, step, serial, device.ReferenceDeviceId, device.ProvinceId, device.CityId, device.DistrictId, device.SchoolId, device.Address, device.FirstPulsePrice, device.SecondPulsePrice, device.ThirdPulsePrice, device.FourthPulsePrice, device.FirstPulseName, device.SecondPulseName, device.ThirdPulseName, device.FourthPulseName, device.Status, device.CreatedAt, device.UpdatedAt)
		sql = sql + val
		if k != len(serialList)-1 {
			sql = sql + ","
		}
	}
	r := transAction.Exec(sql)
	if r.Error != nil {
		common.Logger.Warningln(r.Error)
		transAction.Rollback()
		return r.Error
	}
	transAction.Commit()
	return nil
}

func (self *DeviceService) ListByUserIds(userId ...int) (*[]*model.Device, error) {
	list := []*model.Device{}
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id in (?)", userId).Find(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	return &list, nil
}
func (self *DeviceService) ListByUserId(userId int) (*[]*model.Device, error) {
	list := []*model.Device{}
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id = ?", userId).Find(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	return &list, nil
}
func (self *DeviceService) DailyBill(serialNumber string, billAt string, page int, perPage int) (*[]*map[string]interface{}, error) {
	sql := `
	select value,mobile,
	(case
		when device_mode=1 then 601
		when device_mode=2 then 602
		when device_mode=3 then 603
		when device_mode=4 then 604
		else 0
	end),
	created_at
	from ticket
	where device_serial=? and from_unixtime(created_timestamp,'%Y-%m-%d')=?
	and
	status = 7
	order by created_timestamp desc
	`
	rows, err := common.SodaDB_R.Raw(sql, serialNumber, billAt).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	list := make([]*map[string]interface{}, 0)
	for rows.Next() {
		var address string
		var price float64
		var mobile string
		var washtype int
		var createdAt string

		err := rows.Scan(&price, &mobile, &washtype, &createdAt)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{
			"serialNumber": serialNumber,
			"address":      address,
			"price":        price / 100,
			"mobile":       mobile,
			"washType":     washtype,
			"createdAt":    createdAt,
		}
		list = append(list, &m)
	}

	return &list, nil
}

func (self *DeviceService) BasicByUserId(userId int) (*model.Device, error) {
	device := &model.Device{}
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id = ?", userId).First(device)
	if r.Error != nil {
		return nil, r.Error
	}
	return device, nil
}

func (self *DeviceService) BasicMapByUserId(userId ...int) (map[string]*model.Device, error) {
	list := []*model.Device{}
	deviceMap := make(map[string]*model.Device, 0)
	r := common.SodaMngDB_R.Model(&model.Device{}).Where("user_id in (?)", userId).Find(&list)
	if r.Error != nil {
		return nil, r.Error
	}
	for _, _device := range list {
		deviceMap[_device.SerialNumber] = _device
	}
	return deviceMap, nil
}

func (self *DeviceService) Assign(toUser *model.User, fromUser *model.User, serialNumbers []string) (bool, error) {
	transAction := common.SodaMngDB_WR.Begin()
	assignedAt := time.Now().Local().Format("2006-01-02 15:04:05")
	data := map[string]interface{}{
		"user_id":        toUser.ID,
		"assigned_at":    assignedAt,
		"from_user_id":   fromUser.ID,
		"has_retrofited": 0,
		"status":         0,
	}
	r := transAction.Model(&model.Device{}).Where("serial_number in  (?) ", serialNumbers).Updates(data)
	if r.Error != nil {
		transAction.Rollback()
		common.Logger.Warningln("DB Device Assign:", r.Error.Error())
		return false, r.Error
	}
	transAction.Commit()
	return true, nil
}

func (self *DeviceService) TimedUpdateStatus(status int) (bool, error) {
	var r *gorm.DB
	var err error
	firstTimedDurationStr := viper.GetString("device.firstTimedDuration")
	secondTimedDurationStr := viper.GetString("device.secondTimedDuration")
	thirdTimedDurationStr := viper.GetString("device.thirdTimedDuration")
	fourthTimedDurationStr := viper.GetString("device.fourthTimedDuration")
	firstTimedDuration, err := time.ParseDuration(firstTimedDurationStr)
	if err != nil {
		return false, err
	}
	secondTimedDuration, err := time.ParseDuration(secondTimedDurationStr)
	if err != nil {
		return false, err
	}
	thirdTimedDuration, err := time.ParseDuration(thirdTimedDurationStr)
	if err != nil {
		return false, err
	}
	fourthTimedDuration, err := time.ParseDuration(fourthTimedDurationStr)
	if err != nil {
		return false, err
	}

	switch status {
	case 601:
		_, err = timedUpdateStatus("601", firstTimedDuration)
		break
	case 602:
		_, err = timedUpdateStatus("602", secondTimedDuration)
		break
	case 603:
		_, err = timedUpdateStatus("603", thirdTimedDuration)
		break
	case 604:
		_, err = timedUpdateStatus("604", fourthTimedDuration)
		break
	case 606:
		_, err = timedUpdateStatus("606", secondTimedDuration)
		break
	case 607:
		_, err = timedUpdateStatus("607", thirdTimedDuration)
		break
	case 608:
		_, err = timedUpdateStatus("608", fourthTimedDuration)
		break
	default:
		return false, r.Error
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func timedUpdateStatus(status string, timedDuration time.Duration) (bool, error) {
	timeFormat := "2006-01-02 15:04:05"
	/*r := common.SodaMngDB_WR.Model(&model.Device{}).Where("reference_device_id = 3 and status = "+status+" and updated_at <= ?", time.Now().Add(timedDuration).Format(timeFormat)).Update("status", 2)
	if r.Error != nil {
		common.Logger.Warningln("Timed Update "+status+" Device-STATUS:", r.Error.Error())
		return false, r.Error
	}
	r = common.SodaMngDB_WR.Model(&model.Device{}).Where("reference_device_id = 2 and status = "+status+" and updated_at <= ?", time.Now().Add(timedDuration).Format(timeFormat)).Update("status", 0)
	if r.Error != nil {
		common.Logger.Warningln("Timed Update "+status+" Device-STATUS:", r.Error.Error())
		return false, r.Error
	}*/
	r := common.SodaMngDB_WR.Model(&model.Device{}).Where("status = "+status+" and updated_at <= ?", time.Now().Add(timedDuration).Format(timeFormat)).Update("status", 0)
	if r.Error != nil {
		common.Logger.Warningln("Timed Update "+status+" Device-STATUS:", r.Error.Error())
		return false, r.Error
	}
	return true, nil
}
