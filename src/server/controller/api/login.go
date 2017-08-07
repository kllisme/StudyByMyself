package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	permissionModel "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/payload"
	"maizuo.com/soda/erp/api/src/server/service"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	"strings"
	"github.com/bitly/go-simplejson"
)

type LoginController struct {
}

func (self *LoginController) Login(ctx *iris.Context) {
	var (
		captchaKey = viper.GetString("server.captcha.key")
		userService = service.UserService{}
		tokenService = service.TokenService{}
		menuService = permission.MenuService{}
		userRoleRelService = permission.UserRoleRelService{}
		roleMenuRelService = permission.PermissionMenuRelService{}
		permissionActionRelService = permission.PermissionActionRelService{}
		actionService = permission.ActionService{}
		rolePermissionRelService = permission.RolePermissionRelService{}
		permissionElementRelService = permission.PermissionElementRelService{}
		elementService = permission.ElementService{}
	)

	//每次调用返回时都清一次图片验证码
	defer ctx.Session().Delete(captchaKey)

	params := simplejson.New()
	err := ctx.ReadJSON(&params)
	if err != nil {
		common.Render(ctx, "27010102", nil)
		return
	}
	account := strings.TrimSpace(params.Get("account").MustString())
	password := strings.TrimSpace(params.Get("password").MustString())
	captcha := strings.TrimSpace(params.Get("captcha").MustString())

	/*判断不能为空*/
	if account == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	if password == "" {
		common.Render(ctx, "27010104", nil)
		return
	}
	if captcha == "" {
		common.Render(ctx, "27010105", nil)
		return
	}

	captchaCache := ctx.Session().GetString(captchaKey)
	if captchaCache == "" {
		common.Render(ctx, "27010106", nil)
		return
	}

	if captchaCache != captcha {
		common.Render(ctx, "27010107", nil)
		return
	}

	userEntity, err := userService.GetByAccount(account)
	if err != nil {
		common.Render(ctx, "27010108", nil)
		return
	}
	if userEntity.Password != password {
		common.Render(ctx, "27010109", nil)
		return
	}

	sessionInfo := payload.SessionInfo{
		User:        userEntity,
		MenuList:    &[]*permissionModel.Menu{},
		ActionList:  &[]*permissionModel.Action{},
		ElementList: &[]*permissionModel.Element{},
	}
	//获取权限
	roleIDs, err := userRoleRelService.GetRoleIDsByUserID(userEntity.ID)
	if err != nil {
		common.Render(ctx, "27010112", nil)
		return
	}
	permissionIDs, err := rolePermissionRelService.GetPermissionIDsByRoleIDs(roleIDs)
	if err != nil {
		common.Render(ctx, "27010117", nil)
		return
	}
	if len(permissionIDs) != 0 {
		menuIDs, err := roleMenuRelService.GetMenuIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27010113", nil)
			return
		}
		if len(menuIDs) != 0 {
			menuList, err := menuService.GetListByIDs(menuIDs)
			if err != nil {
				common.Render(ctx, "27010114", nil)
				return
			}
			sessionInfo.MenuList = menuList
		}
		actionIDs, err := permissionActionRelService.GetActionIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27010115", nil)
			return
		}
		if len(actionIDs) != 0 {
			actionList, err := actionService.GetListByIDs(actionIDs)
			if err != nil {
				common.Render(ctx, "27010116", nil)
				return
			}
			sessionInfo.ActionList = actionList
		}
		elementIDs, err := permissionElementRelService.GetElementIDsByPermissionIDs(permissionIDs)
		if err != nil {
			common.Render(ctx, "27010119", nil)
			return
		}
		if len(elementIDs) != 0 {
			elementList, err := elementService.GetListByIDs(elementIDs)
			if err != nil {
				common.Render(ctx, "27010118", nil)
				return
			}
			sessionInfo.ElementList = elementList
		}
	}
	jsonObj, _ := json.Marshal(sessionInfo)
	jsonString := string(jsonObj)
	ctx.Session().Set(viper.GetString("server.session.user.key"), jsonString)
	ctx.Session().Set(viper.GetString("server.session.user.id"), userEntity.ID)
	token, err := tokenService.Token(ctx)
	if err != nil {
		common.Render(ctx, "27010110", err)
		return
	}
	common.Render(ctx, "27010100", token)
}

//微信登录
func (self *LoginController) WechatLogin(ctx *iris.Context) {

}

//手机短信验证码，手机密码登录
func (self *LoginController) PassportLogin(ctx *iris.Context) {

}

func (self *LoginController) Logout(ctx *iris.Context) {
	ctx.SessionDestroy()
	common.Render(ctx, "27010200", nil)
	return
}
