package crm

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service"
	"maizuo.com/soda/erp/api/src/server/model"
	"maizuo.com/soda/erp/api/src/server/common"
)

type DeviceController struct {

}

func (self *DeviceController)Paging(ctx *iris.Context) {
	deviceService := service.DeviceService{}
	userService := service.UserService{}

	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称
	serial := ctx.URLParam("deviceSerial")

	userList := make([]*model.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if _p, err := userService.Paging(keywords, "", 0, 0, 0, 0); err != nil {
			common.Render(ctx, "04020100", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
		if _p, err := userService.Paging("", keywords, 0, 0, 0, 0); err != nil {
			common.Render(ctx, "04020100", err)
			return
		} else {
			_userList := _p.Objects.([]*model.User)
			userList = append(userList, _userList...)
		}
	}
	if len(userList) != 0 {
		for _, user := range userList {
			userIDs = append(userIDs, user.ID)
		}
	}

	pagination, err := deviceService.Paging(userIDs, []int{}, "", serial, 0, 0, []int{}, offset, limit)
	if err != nil {
		common.Render(ctx, "04020100", err)
		return
	}
	referDeviceList, err := deviceService.GetAllReferenceDevice()
	if err != nil {
		common.Render(ctx, "04020100", err)
		return
	}
	deviceList := pagination.Objects.([]*model.Device)
	for _, device := range deviceList {
		user, err := userService.GetById(device.UserID)
		if err != nil {
			common.Render(ctx, "04020401", err)
			return
		}
		fromUser, err := userService.GetById(device.FromUserID)
		if err != nil {
			common.Render(ctx, "04020401", err)
			return
		}
		device.FromUserName = fromUser.Name
		device.FromUserMobile = fromUser.Mobile
		device.UserName = user.Name
		device.UserMobile = user.Mobile
		for _, refer := range *referDeviceList {
			if refer.ID == device.ReferenceDeviceID {
				device.ReferenceDevice = refer.Name
				break
			}
		}
	}

	if err != nil {
		common.Render(ctx, "04020201", nil)
		return
	}
	common.Render(ctx, "04020200", pagination)
	return
}

func (self *DeviceController)Free(ctx *iris.Context) {
	deviceService := service.DeviceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020101", err)
		return
	}
	device, err := deviceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020100", err)
		return
	}
	device.Status = 0
	device, err = deviceService.UpdateStatus(device)
	if err != nil {
		common.Render(ctx, "04020102", err)
	}
	common.Render(ctx, "04020100", device)
}

func (self *DeviceController)Unlock(ctx *iris.Context) {
	deviceService := service.DeviceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020101", err)
		return
	}
	device, err := deviceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020100", err)
		return
	}
	device.Status = 0
	device, err = deviceService.UpdateStatus(device)
	if err != nil {
		common.Render(ctx, "04020102", err)
	}
	common.Render(ctx, "04020100", device)
}

func (self *DeviceController)Lock(ctx *iris.Context) {
	deviceService := service.DeviceService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "04020101", err)
		return
	}
	device, err := deviceService.GetByID(id)
	if err != nil {
		common.Render(ctx, "04020100", err)
		return
	}
	device.Status = 9
	device, err = deviceService.UpdateStatus(device)
	if err != nil {
		common.Render(ctx, "04020102", err)
	}
	common.Render(ctx, "04020100", device)
}

func (self *DeviceController)Remove(ctx *iris.Context) {
	deviceService := service.DeviceService{}
	userService := service.UserService{}
	ids := []int{}
	if err := ctx.ReadJSON(&ids); err != nil {
		common.Render(ctx, "04020401", nil)
		return
	}
	if len(ids) == 0 {
		common.Render(ctx, "04020401", nil)
		return
	}
	deviceList := []*model.Device{}
	for _, id := range ids {
		device, err := deviceService.GetByID(id)
		if err != nil {
			common.Render(ctx, "04020401", err)
			return
		}
		user, err := userService.GetById(device.UserID)
		if err != nil {
			common.Render(ctx, "04020100", err)
			return
		}
		device, err = deviceService.Reset(id, user)
		if err != nil {
			common.Render(ctx, "04020100", err)
			return
		}
		deviceList = append(deviceList, device)
	}

	common.Render(ctx, "04020400", deviceList)
	return
}
