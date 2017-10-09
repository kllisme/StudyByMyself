package api

import (
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/payload"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
)

type UserController struct{}

func (self *UserController) AuthorizationUser(ctx *iris.Context) {

}

func (self *UserController) Paging(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, _ := ctx.URLParamInt("id")
	account := strings.TrimSpace(ctx.URLParam("account"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	roleID, _ := ctx.URLParamInt("roleId")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := userService.Paging(name, account, []int{id}, roleID, offset, limit)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020300", result)
	return
}

func (self *UserController) Create(ctx *iris.Context) {
	userService := mngService.UserService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020201", err)
		return
	}
	account := strings.TrimSpace(params.Get("account").MustString())

	if account == "" {
		common.Render(ctx, "27020202", nil)
		return
	} else if functions.CountRune(account) < 5 || functions.CountRune(account) > 50 {
		common.Render(ctx, "27020208", nil)
		return
	}
	password := strings.TrimSpace(params.Get("password").MustString())
	if password == "" {
		common.Render(ctx, "27020207", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	if name == "" {
		common.Render(ctx, "27020202", nil)
		return
	} else if functions.CountRune(name) > 50 {
		common.Render(ctx, "27020209", nil)
		return
	}
	contact := strings.TrimSpace(params.Get("contact").MustString())
	if contact == "" {
		common.Render(ctx, "27020203", nil)
		return
	}
	mobile := strings.TrimSpace(params.Get("mobile").MustString())
	if mobile == "" {
		common.Render(ctx, "27020204", nil)
		return
	} else if len(mobile) != 11 {
		common.Render(ctx, "27020210", nil)
		return
	}

	telephone := strings.TrimSpace(params.Get("telephone").MustString())

	address := strings.TrimSpace(params.Get("address").MustString())
	if address == "" {
		common.Render(ctx, "27020205", nil)
		return
	}
	parentID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000008", err)
		return
	}

	user := mngModel.User{
		Account:   account,
		Name:      name,
		Mobile:    mobile,
		Contact:   contact,
		ParentID:  parentID,
		Telephone: telephone,
		Address:   address,
		Password:  password,
	}

	if _, err := userService.GetByAccount(account); err == nil {
		common.Render(ctx, "27020206", nil)
		return
	}
	entity, err := userService.Create(&user)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *UserController) Update(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	_, err = userService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	user := mngModel.User{}
	if err := ctx.ReadJSON(&user); err != nil {
		common.Render(ctx, "27020401", err)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		common.Render(ctx, "27020402", nil)
		return
	} else if functions.CountRune(user.Name) > 50 {
		common.Render(ctx, "27020403", nil)
		return
	}
	user.Contact = strings.TrimSpace(user.Contact)
	user.Mobile = strings.TrimSpace(user.Mobile)
	if user.Mobile == "" {
		common.Render(ctx, "27020404", nil)
		return
	} else if len(user.Mobile) != 11 {
		common.Render(ctx, "27020405", nil)
		return
	}
	user.Telephone = strings.TrimSpace(user.Telephone)
	user.Address = strings.TrimSpace(user.Address)
	user.ID = id
	entity, err := userService.Update(&user)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020400", entity)
}

func (self *UserController) Delete(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	if err := userService.DeleteById(id); err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020500", nil)
}

func (self *UserController) AssignRoles(ctx *iris.Context) {
	userService := mngService.UserService{}
	userRoleRelService := mngService.UserRoleRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = userService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	roleIDs := make([]int, 0)
	if err := ctx.ReadJSON(&roleIDs); err != nil {
		common.Render(ctx, "27020901", err)
		return
	}
	result, err := userRoleRelService.AssignRoles(id, roleIDs)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020900", result)
}

func (self *UserController) GetRoles(ctx *iris.Context) {
	userService := mngService.UserService{}
	userRoleRelService := mngService.UserRoleRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	_, err = userService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}

	result, err := userRoleRelService.GetRoleIDsByUserID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27021000", result)
}

func (self *UserController) GetByID(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	userEntity, err := userService.GetByID(id)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020600", userEntity)
}

//GetSessionInfo	use for pull info which shown on pages after login
func (self *UserController) GetProfile(ctx *iris.Context) {
	var (
		userService                 = mngService.UserService{}
		menuService                 = mngService.MenuService{}
		userRoleRelService          = mngService.UserRoleRelService{}
		roleMenuRelService          = mngService.PermissionMenuRelService{}
		rolePermissionRelService    = mngService.RolePermissionRelService{}
		permissionElementRelService = mngService.PermissionElementRelService{}
		elementService              = mngService.ElementService{}
	)

	id, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000008", err)
		return
	}

	userEntity, err := userService.GetByID(id)
	if err != nil {
		common.Render(ctx, "27020120", err)
		return
	}

	sessionInfo := payload.SessionInfo{
		User:        userEntity,
		MenuList:    &[]*mngModel.Menu{},
		ElementList: &[]*mngModel.Element{},
	}
	//获取权限
	roleIDs, err := userRoleRelService.GetRoleIDsByUserID(userEntity.ID)
	if err != nil {
		common.Render(ctx, "27020112", err)
		return
	}
	permissionIDs, err := rolePermissionRelService.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		common.Render(ctx, "27020117", err)
		return
	}
	if len(permissionIDs) != 0 {
		menuIDs, err := roleMenuRelService.GetMenuIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27020113", err)
			return
		}
		if len(menuIDs) != 0 {
			menuList, err := menuService.GetListByIDs(menuIDs)
			if err != nil {
				common.Render(ctx, "27020114", err)
				return
			}
			sessionInfo.MenuList = menuList
		}
		elementIDs, err := permissionElementRelService.GetElementIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27020119", err)
			return
		}
		if len(elementIDs) != 0 {
			elementList, err := elementService.GetListByIDs(elementIDs)
			if err != nil {
				common.Render(ctx, "27020118", err)
				return
			}
			sessionInfo.ElementList = elementList
		}
	}
	//userEntity, _ := userService.GetById(user.ID)
	common.Render(ctx, "27020100", sessionInfo)
}

//ResetPassword 将指定用户的密码重置为服务器默认初始密码
func (self *UserController) ResetPassword(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", err)
		return
	}
	defaultPassword := viper.GetString("defaultPassword")
	user, err := userService.ChangePassword(id, defaultPassword)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020800", user)
}

//ChangePassword 更改当前登录用户的密码
func (self *UserController) ChangePassword(ctx *iris.Context) {
	userService := mngService.UserService{}
	currentUserID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000008", err)
		return
	}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020702", err)
		return
	}

	oldPassword := params.Get("oldPassword").MustString()
	if oldPassword == "" {
		common.Render(ctx, "27020702", nil)
		return
	}
	user := mngModel.User{}
	user.ID = currentUserID
	user.Password = oldPassword
	if _, err := userService.CheckInfo(&user); err != nil {
		common.Render(ctx, "27020701", err)
		return
	}

	newPassword := params.Get("newPassword").MustString()
	if newPassword == "" {
		common.Render(ctx, "27020703", nil)
		return
	}
	entity, err := userService.ChangePassword(currentUserID, newPassword)
	if err != nil {
		common.Render(ctx, "000002", err)
		return
	}
	common.Render(ctx, "27020700", entity)
}
