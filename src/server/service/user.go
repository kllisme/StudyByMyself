package service

import (
	"maizuo.com/soda/cron/src/server/rpc"
	"maizuo.com/soda/erp-api/src/server/common"
)

type UserService struct {
}

func (self *UserService) GetById(id int) (*rpc.User_Out, *common.Result) {
	return nil, nil
}

func (self *UserService) GetByMobile(mobile string) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) GetList(limit int, offset int) (*rpc.User_GetList_Out, error) {

	return nil, nil
}

func (self *UserService) Create(in *rpc.User_Create_In) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) UpdateById(in *rpc.User_UpdateById_In) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) UpdateByMobile(in *rpc.User_UpdateByMobile_In) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) DeleteById(id int) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) DeleteByMobile(mobile string) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) RemoveById(id int) (*rpc.User_Out, error) {
	return nil, nil
}

func (self *UserService) RemoveByMobile(mobile string) (*rpc.User_Out, error) {
	return nil, nil
}
