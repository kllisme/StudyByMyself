package two

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/two"
)

type ChannelService struct {

}

func (self *ChannelService)GetByID(id int) (*two.Channel, error) {
	channel := two.Channel{}
	err := common.Soda2DB_R.Where(id).Find(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (self *ChannelService)GetAll() (*[]*two.Channel, error) {
	channelList := make([]*two.Channel, 0)
	if err := common.Soda2DB_R.Order("id desc").Find(&channelList).Error; err != nil {
		return nil, err
	}
	return &channelList, nil

}
