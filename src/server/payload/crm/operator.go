package crm

import (
	"time"
	model "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	service "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	"github.com/jinzhu/gorm"
	"errors"
	"maizuo.com/soda/erp/api/src/server/common"
)

type Operator struct {
	ID                   int        `json:"id"`
	Name                 string        `json:"name"`
	Role                 string        `json:"role"`
	Account              string        `json:"account"`
	CreatedAt            time.Time        `json:"createdAt"`
	Contact              string        `json:"contact"`
	Mobile               string        `json:"mobile"`
	Telephone            string        `json:"telephone"`
	Address              string        `json:"address"`
	ParentOperator       string        `json:"parentOperator"`
	ParentOperatorMobile string        `json:"parentOperatorMobile"`
	DeviceCount          int        `json:"deviceCount"`
	AutoSettle           string        `json:"autoSettle"`
	Payment              string        `json:"payment"`
	CashAccount          string        `json:"cashAccount"`
}

func (self *Operator) Map(user *model.User) error {
	userService := service.UserService{}
	userCashAccountService := service.UserCashAccountService{}
	deviceService := service.DeviceService{}
	roleService := service.RoleService{}
	userRoleService := service.UserRoleRelService{}
	self.ID = user.ID
	self.Name = user.Name

	roleIDs, err := userRoleService.GetRoleIDsByUserID(user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			self.Role = "无"
		} else {
			common.Logger.Error("用户角色信息出错")
			return errors.New("用户角色信息出错")
		}
	} else {
		for _, roleID := range roleIDs {
			_role, err := roleService.GetByID(roleID)
			if err != nil {
				return errors.New("角色信息出错")
			}
			self.Role += _role.Name + " "
		}
	}
	self.Account = user.Account
	self.CreatedAt = user.CreatedAt
	self.Contact = user.Contact
	self.Mobile = user.Mobile
	self.Telephone = user.Telephone
	self.Address = user.Address

	parent, err := userService.GetByID(user.ParentID)
	if err != nil {
		return errors.New("用户父级信息出错")
	}
	self.ParentOperator = parent.Name
	self.ParentOperatorMobile = parent.Mobile
	self.DeviceCount, err = deviceService.TotalByUser(user.ID)
	if err != nil {
		return errors.New("设备信息出错")
	}
	cashAccount, err := userCashAccountService.GetByUserID(user.ID)
	if err != nil {
		return errors.New("收款账号信息出错")
	}
	switch cashAccount.Mode {
	case 0:
		self.AutoSettle = "是"
	case 1:
		self.AutoSettle = "否"
	}
	switch cashAccount.Type {
	case 1:
		self.Payment = "支付宝"
		self.CashAccount = cashAccount.RealName+"|"+cashAccount.Account
	case 2:
		self.Payment = "微信"
		self.CashAccount = cashAccount.RealName

	case 3:
		self.Payment = "银行"
		self.CashAccount = cashAccount.RealName
	default:
		self.Payment = "无"
		self.CashAccount = cashAccount.RealName
	}
	return nil
}
