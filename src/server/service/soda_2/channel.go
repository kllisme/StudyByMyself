package soda_2

import (
	"maizuo.com/soda/erp/api/src/server/common"
	twoModel "maizuo.com/soda/erp/api/src/server/model/soda_2"
)

type ChannelService struct {

}

func (self *ChannelService)GetByID(id int) (*twoModel.Channel, error) {
	channel := twoModel.Channel{}
	err := common.Soda2DB_R.Where(id).Find(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (self *ChannelService)GetAll() (*[]*twoModel.Channel, error) {
	channelList := make([]*twoModel.Channel, 0)
	if err := common.Soda2DB_R.Order("id desc").Find(&channelList).Error; err != nil {
		return nil, err
	}
	return &channelList, nil

}
