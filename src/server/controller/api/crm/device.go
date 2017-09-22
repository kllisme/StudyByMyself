package crm

import (
	"gopkg.in/kataras/iris.v5"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/spf13/viper"
)

type DeviceController struct {

}

func (self *DeviceController)Paging(ctx *iris.Context) {
	deviceService := mngService.DeviceService{}
	userService := mngService.UserService{}
	deviceOperateService := &mngService.DeviceOperateService{}
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称
	serial := ctx.URLParam("deviceSerial")

	userList := make([]*mngModel.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if _p, err := userService.Paging(keywords, "", 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05020101", err)
			return
		} else {
			_userList := _p.Objects.([]*mngModel.User)
			userList = append(userList, _userList...)
		}
		if _p, err := userService.Paging("", keywords, 0, 0, 0, 0); err != nil {
			common.Render(ctx, "05020102", err)
			return
		} else {
			_userList := _p.Objects.([]*mngModel.User)
			userList = append(userList, _userList...)
		}
		if len(userList) != 0 {
			for _, user := range userList {
				userIDs = append(userIDs, user.ID)
			}
		} else {
			userIDs = []int{-1}
		}
	}

	pagination, err := deviceService.Paging(userIDs, []int{}, "", serial, 0, 0, []int{}, offset, limit)
	if err != nil {
		common.Render(ctx, "05020103", err)
		return
	}
	referDeviceList, err := deviceService.GetAllReferenceDevice()
	if err != nil {
		common.Render(ctx, "05020104", err)
		return
	}
	deviceList := pagination.Objects.([]*mngModel.Device)
	for _, device := range deviceList {
		user, err := userService.GetById(device.UserID)
		if err != nil {
			common.Render(ctx, "05020105", err)
			return
		}
		operationList, err := deviceOperateService.GetBySerialNumber(device.SerialNumber)
		if err != nil {
			common.Render(ctx, "05020106", err)
			return
		}
		for _,operation := range *operationList {
			if operation.OperatorType == 1 || operation.OperatorType == 4 {
				device.AssignerID = operation.OperatorID
				break
			}
		}
		if device.AssignerID != 0 {
			assigner, err := userService.GetById(device.AssignerID)
			if err != nil {
				common.Render(ctx, "05020107", err)
				return
			}
			device.Assigner = assigner.Name
			device.AssignerMobile = assigner.Mobile
		}
		device.UserName = user.Name
		device.UserMobile = user.Mobile
		for _, refer := range *referDeviceList {
			if refer.ID == device.ReferenceDeviceID {
				device.ReferenceDevice = refer.Name
				break
			}
		}
	}
	common.Render(ctx, "05020100", pagination)
	return
}

func (self *DeviceController)UpdateStatus(ctx *iris.Context) {
	deviceService := mngService.DeviceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "05020201", err)
		return
	}
	device := mngModel.Device{}
	if err := ctx.ReadJSON(&device); err != nil {
		common.Render(ctx, "05020202", err)
		return
	}
	if device.Status != 0 && device.Status != 9 {
		common.Render(ctx, "05020203", err)
		return
	}

	_device, err := deviceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "05020204", err)
		return
	}
	device.ID = id

	if _device, err = deviceService.UpdateStatus(&device); err != nil {
		common.Render(ctx, "05020205", err)
		return
	}

	if device.Status == 0 {
		deviceService.UnLockDevice(_device.SerialNumber)
	}
	common.Render(ctx, "05020200", _device)
}

func (self *DeviceController)Reset(ctx *iris.Context) {
	deviceService := mngService.DeviceService{}
	userService := mngService.UserService{}
	deviceOperateService := &mngService.DeviceOperateService{}
	operatorId, _ := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	ids := []int{}
	if err := ctx.ReadJSON(&ids); err != nil {
		common.Render(ctx, "05020301", err)
		return
	}
	if len(ids) == 0 {
		common.Render(ctx, "05020302", nil)
		return
	}
	deviceList := []*mngModel.Device{}
	for _, id := range ids {
		device, err := deviceService.GetByID(id)
		if err != nil {
			common.Render(ctx, "05020303", err)
			return
		}
		user, err := userService.GetById(device.UserID)
		if err != nil {
			common.Render(ctx, "05020304", err)
			return
		}
		device, err = deviceService.Reset(id, user)
		if err != nil {
			common.Render(ctx, "05020305", err)
			return
		}
		deviceOperation := &mngModel.DeviceOperate{
			OperatorID:   operatorId,
			OperatorType: 1,
			SerialNumber: device.SerialNumber,
			UserID:       device.UserID,
			FromUserID:   device.FromUserID,
			ToUserID:     1,
		}
		_, err = deviceOperateService.Create(deviceOperation)
		if err != nil {
			common.Render(ctx, "05020306", err)
			return
		}
		deviceList = append(deviceList, device)
	}

	common.Render(ctx, "05020300", deviceList)
	return
}
