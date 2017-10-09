package crm

import (
	"github.com/bitly/go-simplejson"
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/jinzhu/gorm"
	mngService "maizuo.com/soda/erp/api/src/server/service/soda_manager"
	"maizuo.com/soda/erp/api/src/server/payload/crm"
	"strings"
	mngModel "maizuo.com/soda/erp/api/src/server/model/soda_manager"
	"maizuo.com/soda/erp/api/src/server/kit/functions"
)

type OperatorController struct{}

func (self *OperatorController) GetByID(ctx *iris.Context) {
	userService := mngService.UserService{}

	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "05050101", nil)
		return
	}
	operator := crm.Operator{}
	userEntity, err := userService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			common.Render(ctx, "05050102", err)
			return
		}
		common.Render(ctx, "05050103", err)
		return
	}
	err = operator.Map(userEntity)
	if err != nil {
		common.Render(ctx, "05050104", err)
		return
	}
	common.Render(ctx, "05050100", operator)
}

func (self *OperatorController) ChangePassword(ctx *iris.Context) {
	userService := mngService.UserService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "05050201", err)
		return
	}
	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "05050202", err)
		return
	}
	newPassword := strings.TrimSpace(params.Get("newPassword").MustString())
	if len(newPassword) < 6 || len(newPassword) > 16 {
		common.Render(ctx, "05050203", nil)
		return
	}
	entity, err := userService.ChangePassword(id, newPassword)
	if err != nil {
		common.Render(ctx, "05050204", err)
		return
	}
	common.Render(ctx, "05050200", entity)
}

func (self *OperatorController) Paging(ctx *iris.Context) {
	userService := mngService.UserService{}
	limit, _ := ctx.URLParamInt("limit")       // Default: 10
	offset, _ := ctx.URLParamInt("offset")     //  Default: 0 列表起始位:
	keywords := ctx.URLParam("keywords")               // 运营商名称、帐号名称

	_userList := make([]*mngModel.User, 0)
	userIDs := make([]int, 0)
	if keywords != "" {
		if __userList, err := userService.GetList(keywords, "", []int{}, 0); err != nil {
			common.Render(ctx, "05050301", err)
			return
		} else {
			_userList = append(_userList, *__userList...)
		}
		if __userList, err := userService.GetList("", keywords, []int{}, 0); err != nil {
			common.Render(ctx, "05050302", err)
			return
		} else {
			_userList = append(_userList, *__userList...)
		}
		if len(_userList) != 0 {
			for _, user := range _userList {
				userIDs = append(userIDs, user.ID)
			}
		} else {
			userIDs = []int{-1}
		}
		userIDs = functions.Uniq(userIDs)
	}

	pagination, err := userService.Paging("", "", userIDs, 0, offset, limit)
	if err != nil {
		common.Render(ctx, "05050303", err)
		return
	}
	if pagination.Pagination.Total != 0 {
		operatorList := make([]*crm.Operator, 0)
		userList := pagination.Objects.([]*mngModel.User)
		for _, user := range userList {
			operator := crm.Operator{}
			err = operator.Map(user)
			if err != nil {
				common.Render(ctx, "05050304", err)
				return
			}
			operatorList = append(operatorList, &operator)
		}
		pagination.Objects = operatorList

	}
	common.Render(ctx, "05050300", pagination)
	return
}
