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

	common.Logger.Debug(name, account, id, roleID, page, perPage)
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
	if account == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	//name := strings.TrimSpace(params.MustString("name"))
	contact := strings.TrimSpace(params.Get("contact").MustString())
	if contact == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	//mobile := strings.TrimSpace(params.MustString("mobile"))
	parentID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "000001", nil)
		return
	}
	telephone := strings.TrimSpace(params.Get("telephone").MustString())
	if telephone == "" {
		common.Render(ctx, "27010103", nil)
		return
	}

	address := strings.TrimSpace(params.Get("address").MustString())
	if address == "" {
		common.Render(ctx, "27010103", nil)
		return
	}
	user := model.User{
		Account:account,
		Contact:contact,
		ParentID:parentID,
		Telephone:telephone,
		Address:address,
	}
	entity, err := userService.Create(&user)
	if err != nil {
		common.Render(ctx, "27010103", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *UserController)Update(ctx *iris.Context) {
	//TODO Complete status codes
	userService := service.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	user, err := userService.GetById(id)
	if err != nil {
		common.Render(ctx, "000003", nil)
		return
	}

	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "27020401", nil)
		return
	}
	account := strings.TrimSpace(params.Get("account").MustString())
	if account == "" {
		common.Render(ctx, "", nil)
		return
	}
	name := strings.TrimSpace(params.Get("name").MustString())
	contact := strings.TrimSpace(params.Get("contact").MustString())
	if contact == "" {
		common.Render(ctx, "", nil)
		return
	}
	mobile := strings.TrimSpace(params.Get("mobile").MustString())
	parentID := params.Get("parentId").MustInt()
	telephone := strings.TrimSpace(params.Get("telephone").MustString())
	if telephone == "" {
		common.Render(ctx, "", nil)
		return
	}
	address := strings.TrimSpace(params.Get("address").MustString())
	if address == "" {
		common.Render(ctx, "", nil)
		return
	}
	user.Account = account
	user.Name = name
	user.Contact = contact
	user.ParentID = parentID
	user.Telephone = telephone
	user.Address = address
	user.Mobile = mobile
	entity, err := userService.UpdateById(user)
	if err != nil {
		common.Render(ctx, "27010103", nil)
		return
	}
	common.Render(ctx, "27010100", entity)
}

func (self *UserController)Delete(ctx *iris.Context) {
	//userService := service.UserService{}
}

func (self *UserController)AssignRoles(ctx *iris.Context) {

}

func (self *UserController) GetSessionUser(ctx *iris.Context) {
	currentUserID, err := ctx.Session().GetInt(viper.GetString("server.session.user.id"))
	if err != nil {
		common.Render(ctx, "", nil)
	}
	userService := service.UserService{}
	userEntity, _ := userService.GetById(currentUserID)
	common.Render(ctx, "", userEntity)
}

func (self *UserController) GetById(ctx *iris.Context) {
	userService := service.UserService{}
	id, err := ctx.URLParamInt("id")
	if err != nil {
		common.Render(ctx, "", nil)
	}
	userEntity, err := userService.GetById(id)
	if err != nil {
		common.Render(ctx, "", nil)
	}
	common.Render(ctx, "", userEntity)
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

}
