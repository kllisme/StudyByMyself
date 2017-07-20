package api

import (
	"github.com/bitly/go-simplejson"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	"github.com/spf13/viper"
	"encoding/json"
	"maizuo.com/soda/erp/api/src/server/service/permission"
	permissionModel        "maizuo.com/soda/erp/api/src/server/model/permission"
	"maizuo.com/soda/erp/api/src/server/payload"
	"strings"
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
		roleMenuRelService = permission.RoleMenuRelService{}
		roleAPIRelService = permission.RoleAPIRelService{}
		apiService = permission.APIService{}
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
		User:userEntity,
		MenuList:&[]*permissionModel.Menu{},
		APIList:&[]*permissionModel.API{},
	}
	//获取权限
	roleIDs, err := userRoleRelService.GetRoleIDsByUserID(userEntity.Id)
	if err != nil {
		common.Render(ctx, "27010112", nil)
		return
	}
	if len(roleIDs) != 0 {
		menuIDs, err := roleMenuRelService.GetMenuIDsByRoleIDs(roleIDs)
		if err != nil {
			common.Render(ctx, "27010113", nil)
			return
		}
		if len(menuIDs) != 0 {
			menuList, err := menuService.GetListByIds(menuIDs)
			if err != nil {
				common.Render(ctx, "27010114", nil)
				return
			}
			sessionInfo.MenuList = menuList
		}
		apiIDs, err := roleAPIRelService.GetAPIIDsByRoleIDs(roleIDs)
		if err != nil {
			common.Render(ctx, "27010115", nil)
			return
		}
		if len(apiIDs) != 0 {
			apiList, err := apiService.GetListByIds(apiIDs)
			if err != nil {
				common.Render(ctx, "27010116", nil)
				return
			}
			sessionInfo.APIList = apiList
		}
	}
	jsonObj, _ := json.Marshal(sessionInfo)
	jsonString := string(jsonObj)
	ctx.Session().Set(viper.GetString("server.session.user.key"), jsonString)

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
	common.Render(ctx, "27010111", nil)
	return
}
