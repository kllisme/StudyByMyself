package two

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api/two"
)

func Setup(v iris.MuxAPI) {
	var (
		topicController=two.TopicController{}
		channelController=two.ChannelController{}

	)
	_api := v.Party("/2")

	_api.Get("/circles",topicController.PagingCircle)

	_api.Get("/circles/summary", topicController.GetSummary)

	_api.Get("/topics", topicController.Paging)

	_api.Get("/topics/:id", topicController.GetByID)

	_api.Put("/topics/:id/channel", topicController.UpdateChannel)

	_api.Put("/topics/:id/status", topicController.UpdateStatus)



	_api.Get("/channels", channelController.GetAll)
}
