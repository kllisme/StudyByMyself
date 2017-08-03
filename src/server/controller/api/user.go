package api

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service"
	"github.com/spf13/viper"
	"github.com/square/go-jose/json"
	"maizuo.com/soda/erp/api/src/server/payload"
	"github.com/bitly/go-simplejson"
	"maizuo.com/soda/erp/api/src/server/model"
	"strings"
	"maizuo.com/soda/erp/api/src/server/service/permission"
)

type UserController struct{}

func (self *UserController) AuthorizationUser(ctx *iris.Context) {

}

func (self *UserController)Paging(ctx *iris.Context) {
	userService := service.UserService{}
	id, _ := ctx.URLParamInt("id")
	account := strings.TrimSpace(ctx.URLParam("account"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	roleID, _ := ctx.URLParamInt("role_id")
	page, _ := ctx.URLParamInt("page")
	perPage, _ := ctx.URLParamInt("per_page")
	result, err := userService.Paging(name, account, id, roleID, page, perPage)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020300", result)
	return
}

func (self *UserController)Create(ctx *iris.Context) {
	userService := service.UserService{}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020201", nil)
		return
	}
	account := strings.TrimSpace(params.Get("account").MustString())
	if _, err := userService.GetByAccount(account); err == nil {
		common.Render(ctx, "27020206", nil)
		return
	}
	if account == "" {
		common.Render(ctx, "27020202", nil)
		return
	}
	password := strings.TrimSpace(params.Get("password").MustString())
	if password == "" {
		common.Render(ctx, "27020207", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	contact := strings.TrimSpace(params.Get("contact").MustString())
	if contact == "" {
		common.Render(ctx, "27020203", nil)
		return
	}
	mobile := strings.TrimSpace(params.Get("mobile").MustString())
	parentID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000001", nil)
		return
	}
	telephone := strings.TrimSpace(params.Get("telephone").MustString())
	if telephone == "" {
		common.Render(ctx, "27020204", nil)
		return
	}

	address := strings.TrimSpace(params.Get("address").MustString())
	if address == "" {
		common.Render(ctx, "27020205", nil)
		return
	}
	user := model.User{
		Account:account,
		Name:name,
		Mobile:mobile,
		Contact:contact,
		ParentID:parentID,
		Telephone:telephone,
		Address:address,
		Password:password,
	}
	entity, err := userService.Create(&user)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *UserController)Update(ctx *iris.Context) {
	userService := service.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	_, err = userService.GetById(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	user := model.User{}
	if err := ctx.ReadJSON(&user); err != nil {
		common.Render(ctx, "27020401", nil)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Contact = strings.TrimSpace(user.Contact)
	user.Mobile = strings.TrimSpace(user.Mobile)
	user.Telephone = strings.TrimSpace(user.Telephone)
	user.Address = strings.TrimSpace(user.Address)
	user.Name = strings.TrimSpace(user.Name)
	user.ID = id
	entity, err := userService.Update(&user)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020400", entity)
}

func (self *UserController)Delete(ctx *iris.Context) {
	userService := service.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	if err := userService.DeleteById(id); err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020500", nil)
}

func (self *UserController)AssignRoles(ctx *iris.Context) {
	userService := service.UserService{}
	userRoleRelService := permission.UserRoleRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = userService.GetById(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	roleIDs := make([]int, 0)
	if err := ctx.ReadJSON(&roleIDs); err != nil {
		common.Render(ctx, "27020901", nil)
		return
	}
	result, err := userRoleRelService.AssignRoles(id, roleIDs)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020900", result)
}

func (self *UserController)GetRoles(ctx *iris.Context) {
	userService := service.UserService{}
	userRoleRelService := permission.UserRoleRelService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	_, err = userService.GetById(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	result, err := userRoleRelService.GetRoleIDsByUserID(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27021000", result)
}

func (self *UserController) GetById(ctx *iris.Context) {
	userService := service.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	userEntity, err := userService.GetById(id)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020600", userEntity)
}

//GetSessionInfo	use for pull info which shown on pages after login
func (self *UserController)GetSessionInfo(ctx *iris.Context) {
	sessionInfo := new(payload.SessionInfo)
	currentUserJson := ctx.Session().GetString(viper.GetString("server.session.user.key"))
	if err := json.Unmarshal([]byte(currentUserJson), sessionInfo); err != nil {
		common.Render(ctx, "27020101", nil)
	}
	//userEntity, _ := userService.GetById(user.ID)
	common.Render(ctx, "27020100", sessionInfo)
}

func (self *UserController)ResetPassword(ctx *iris.Context) {
	userService := service.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}
	defaultPassword := viper.GetString("defaultPassword")
	user, err := userService.ChangePassword(id, defaultPassword)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020800", user)
}

func (self *UserController)ChangePassword(ctx *iris.Context) {
	userService := service.UserService{}
	currentUserID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000001", nil)
		return
	}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "", nil)
		return
	}

	oldPassword := params.Get("oldPassword").MustString()
	if oldPassword == "" {
		common.Render(ctx, "27020702", nil)
		return
	}
	user := model.User{}
	user.ID = currentUserID
	user.Password = oldPassword
	if _, err := userService.CheckInfo(&user); err != nil {
		common.Render(ctx, "27020701", nil)
		return
	}

	newPassword := params.Get("newPassword").MustString()
	if newPassword == "" {
		common.Render(ctx, "27020703", nil)
		return
	}
	entity, err := userService.ChangePassword(currentUserID, newPassword)
	if err != nil {
		common.Render(ctx, "000002", nil)
		return
	}
	common.Render(ctx, "27020700", entity)
}
