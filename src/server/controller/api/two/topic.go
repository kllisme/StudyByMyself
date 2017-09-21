package two

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/service/two"
	"maizuo.com/soda/erp/api/src/server/common"
	"strings"
	"github.com/bitly/go-simplejson"
	model "maizuo.com/soda/erp/api/src/server/model/two"
)

type TopicController struct {

}

func (self *TopicController)PagingCircle(ctx *iris.Context) {
	//CircleList := make([]*payload.Circle,0)
	topicService := two.TopicService{}

	provinceID, _ := ctx.URLParamInt("provinceId")

	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	result, err := topicService.PagingCircle(offset, limit, provinceID)

	if err != nil {
		common.Render(ctx, "03010101", err)
		return
	}
	common.Render(ctx, "03010100", result)
}

func (self *TopicController)GetByID(ctx *iris.Context) {
	topicService := two.TopicService{}
	userService := two.UserService{}
	inboxService:=two.InboxService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "03010301", err)
		return
	}
	topic, err := topicService.GetByID(id)
	if err != nil {
		common.Render(ctx, "03010302", err)
		return
	}
	user, err := userService.GetByID(topic.UserID)
	if err != nil {
		topic.UserName = "神秘用户"
	}
	topic.UserName = user.Name
	consultation, err := inboxService.CountConsultation(id,topic.UserID)
	if err != nil {
		common.Render(ctx, "03010303", err)
		return
	}
	topic.Consultation = consultation
	common.Render(ctx, "03010300", topic)
}

func (self *TopicController)GetSummary(ctx *iris.Context) {
	result := map[string]interface{}{}
	topicService := two.TopicService{}
	userService := two.UserService{}

	topicsCount, err := topicService.CountByCityIDs()
	if err != nil {
		common.Render(ctx, "03010601", err)
		return
	}
	usersCount, err := userService.Count(0)
	if err != nil {
		common.Render(ctx, "03010602", err)
		return
	}
	//圈子数的计算方法为所有有商品发布过的城市
	circleCount, err := topicService.CountCities()
	if err != nil {
		common.Render(ctx, "03010603", err)
		return
	}
	result["topicsCount"] = topicsCount
	result["usersCount"] = usersCount
	result["circleCount"] = circleCount

	common.Render(ctx, "03010600", result)
}

func (self *TopicController)Paging(ctx *iris.Context) {
	userIDs := make([]int, 0)
	userService := two.UserService{}
	topicService := two.TopicService{}
	keywords := strings.TrimSpace(ctx.URLParam("keywords"))
	name := strings.TrimSpace(ctx.URLParam("name"))
	schoolName := strings.TrimSpace(ctx.URLParam("schoolName"))
	statusStr := ctx.URLParam("status")
	status := -1
	if statusStr != "" {
		status, _ = ctx.URLParamInt("status")
	}
	cityID, _ := ctx.URLParamInt("cityId")
	channelID, _ := ctx.URLParamInt("channelId")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	if name != "" {
		result, err := userService.FilterUserIDs(name)
		if err != nil {
			common.Render(ctx, "03010201", err)
			return
		}
		if len(result) == 0 {
			userIDs = append(userIDs,0)
		} else {
			userIDs = result
		}

	}
	result, err := topicService.Paging(cityID, keywords, schoolName, channelID, status, offset, limit, userIDs)
	if err != nil {
		common.Render(ctx, "03010202", err)
		return
	}
	topicList := result.Objects.([]*model.Topic)
	for _, topic := range topicList {
		user, err := userService.GetByID(topic.UserID)
		if err != nil {
			topic.UserName = "神秘用户"
		}
		topic.UserName = user.Name
	}
	common.Render(ctx, "03010200", result)
}

//func (self *TopicController)Delete(ctx *iris.Context) {
//	topicService := two.TopicService{}
//	id, err := ctx.ParamInt("id")
//	if err != nil {
//		common.Render(ctx, "000003", err)
//	}
//	if err := topicService.Delete(id); err != nil {
//		common.Render(ctx, "000002", err)
//	}
//	common.Render(ctx, "27040400", nil)
//}

//func (self *TopicController)Update(ctx *iris.Context) {
//	topicService := two.TopicService{}
//	topic := model.Topic{}
//	id, err := ctx.ParamInt("id")
//	if err != nil {
//		common.Render(ctx, "000003", nil)
//	}
//	_, err = topicService.GetByID(id)
//	if err != nil {
//		common.Render(ctx, "000003", nil)
//	}
//	err = ctx.ReadJSON(&topic)
//	if err != nil {
//		common.Render(ctx, "27040501", nil)
//		return
//	}
//
//	topic.API = strings.TrimSpace(topic.API)
//	topic.Description = strings.TrimSpace(topic.Description)
//	topic.HandlerName = strings.TrimSpace(topic.HandlerName)
//	topic.Method = strings.TrimSpace(topic.Method)
//	topic.ID = id
//	result, err := topicService.Update(&topic)
//	if err != nil {
//		common.Render(ctx, "000002", nil)
//	}
//	common.Render(ctx, "27040500", result)
//}

func (self *TopicController)UpdateChannel(ctx *iris.Context) {
	topicService := two.TopicService{}
	channelService := two.ChannelService{}
	//topic := model.Topic{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "03010401", err)
	}
	topic, err := topicService.GetByID(id)

	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "03010402", err)
		return
	}

	json, e := params.CheckGet("channelId")
	if !e {
		common.Render(ctx, "03010403", nil)
		return
	}
	channelID := json.MustInt()
	channel, err := channelService.GetByID(channelID)
	if err != nil {
		common.Render(ctx, "03010404", err)
		return
	}
	topic.ID = id
	topic.ChannelID = channelID
	topic.ChannelTitle = channel.Title
	result, err := topicService.UpdateChannel(topic)
	if err != nil {
		common.Render(ctx, "03010405", err)
	}
	common.Render(ctx, "03010400", result)
}

func (self *TopicController)UpdateStatus(ctx *iris.Context) {
	topicService := two.TopicService{}
	//topic := model.Topic{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "03010501", err)
	}
	topic, err := topicService.GetByID(id)

	params := simplejson.New()
	if err := ctx.ReadJSON(&params); err != nil {
		common.Render(ctx, "03010502", err)
		return
	}
	status, e := params.CheckGet("status")
	if !e {
		common.Render(ctx, "03010503", nil)
		return
	}
	topic.Status = status.MustInt()
	result, err := topicService.UpdateStatus(topic)
	if err != nil {
		common.Render(ctx, "03010505", err)
	}
	common.Render(ctx, "03010500", result)
}
