package service

import (
	conf "MyFirstProject/config"
	"MyFirstProject/consts"
	"MyFirstProject/pkg/utils/ctl"
	"MyFirstProject/pkg/utils/log"
	util "MyFirstProject/pkg/utils/upload"
	"MyFirstProject/repository/db/dao"
	"MyFirstProject/repository/db/model"
	"MyFirstProject/types"
	"context"
	"mime/multipart"
	"strconv"
	"sync"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once

type ProductSrv struct {
}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{}
	})
	return ProductSrvIns
}

//创建商品

func (s *ProductSrv) ProductCreate(ctx context.Context, files []*multipart.FileHeader, req *types.ProductCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Infoln(err)
		return nil, err
	}
	uId := u.Id
	boss, _ := dao.NewUserDao(ctx).GetUserById(uId)
	// 以第一张作为封面图
	tmp, _ := files[0].Open()
	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name)
	} else {
		path, err = util.UploadToQiNiu(tmp, files[0].Size)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       path,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        true,
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		num := strconv.Itoa(index)
		tmp, _ = file.Open()
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name+num)
		} else {
			path, err = util.UploadToQiNiu(tmp, file.Size)
		}
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		err = dao.NewProductImgDaoByDB(productDao.DB).CreateProductImg(productImg)
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		wg.Done()
	}

	wg.Wait()

	return
}
