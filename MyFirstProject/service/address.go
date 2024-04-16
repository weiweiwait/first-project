package service

import (
	"MyFirstProject/pkg/utils/ctl"
	util "MyFirstProject/pkg/utils/log"
	"MyFirstProject/repository/db/dao"
	"MyFirstProject/repository/db/model"
	"MyFirstProject/types"
	"context"
	"sync"
)

var AddressSrvIns *AddressSrv
var AddressSrvOnce sync.Once

type AddressSrv struct {
}

func GetAddressSrv() *AddressSrv {
	AddressSrvOnce.Do(func() {
		AddressSrvIns = &AddressSrv{}
	})
	return AddressSrvIns
}

//增加地址

func (s *AddressSrv) AddressCreate(ctx context.Context, req *types.AddressCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	addressDao := dao.NewAddressDao(ctx)
	address := &model.Address{
		UserID:  u.Id,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	err = addressDao.CreateAddress(address)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	return
}

//展示某个地址

func (s *AddressSrv) AddressShow(ctx context.Context, req *types.AddressGetReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	address, err := dao.NewAddressDao(ctx).GetAddressByAid(req.Id, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	resp = &types.AddressResp{
		ID:        address.ID,
		UserID:    address.UserID,
		Name:      address.Name,
		Phone:     address.Phone,
		Address:   address.Address,
		CreatedAt: address.CreatedAt.Unix(),
	}

	return
}

//展示所有地址

func (s *AddressSrv) AddressList(ctx context.Context, req *types.AddressListReq) (resp interface{}, err error) {
	u, _ := ctl.GetUserInfo(ctx)
	resp, err = dao.NewAddressDao(ctx).
		ListAddressByUid(u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	return
}

//修改收货地址

func (s *AddressSrv) AddressUpdate(ctx context.Context, req *types.AddressServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	addressDao := dao.NewAddressDao(ctx)
	address := &model.Address{
		UserID:  u.Id,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
	}
	err = addressDao.UpdateAddressById(req.Id, address)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	var addresses []*types.AddressResp
	addresses, err = addressDao.ListAddressByUid(u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	resp = &types.DataListResp{
		Item:  addresses,
		Total: int64(len(addresses)),
	}

	return
}
