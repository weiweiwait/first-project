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
