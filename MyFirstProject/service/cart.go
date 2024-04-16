package service

import (
	"MyFirstProject/pkg/utils/ctl"
	util "MyFirstProject/pkg/utils/log"
	"MyFirstProject/repository/db/dao"
	"MyFirstProject/types"
	"context"
	"sync"
)

var CartSrvIns *CartSrv
var CartSrvOnce sync.Once

type CartSrv struct {
}

func GetCartSrv() *CartSrv {
	CartSrvOnce.Do(func() {
		CartSrvIns = &CartSrv{}
	})
	return CartSrvIns
}

//创建购物车

func (s *CartSrv) CartCreate(ctx context.Context, req *types.CartCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	// 判断有无这个商品
	_, err = dao.NewProductDao(ctx).GetProductById(req.ProductId)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
}
