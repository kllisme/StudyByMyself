package two

import (
	"gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/service/two"
)

type ChannelController struct {

}

func (self *ChannelController)GetByID(ctx *iris.Context) {
	topicService := two.TopicService{}
	id, err := ctx.ParamInt("id")
	if err != nil {
		common.Render(ctx, "03020102", err)
		return
	}
	topic, err := topicService.GetByID(id)
	if err != nil {
		common.Render(ctx, "03020101", err)
		return
	}
	common.Render(ctx, "03020100", topic)
}

func (self *ChannelController)GetAll(ctx *iris.Context) {
	channelService := two.ChannelService{}
	channelList, err := channelService.GetAll()
	if err != nil {
		common.Render(ctx, "03020201", nil)
		return
	}
	common.Render(ctx, "03020200", channelList)
	return
}
