package community

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"maizuo.com/soda/erp/api/src/server/model/community"
)

type ChannelService struct {

}

func (self *ChannelService)GetByID(id int) (*community.Channel, error) {
	channel := community.Channel{}
	err := common.SodaDB_R.Where(id).Find(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (self *ChannelService)GetAll() (*[]*community.Channel, error) {
	channelList := make([]*community.Channel, 0)
	if err := common.SodaDB_R.Order("id desc").Find(&channelList).Error; err != nil {
		return nil, err
	}
	return &channelList, nil

}
