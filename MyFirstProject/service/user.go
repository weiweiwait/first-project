package service

import (
	conf "MyFirstProject/config"
	"MyFirstProject/consts"
	"MyFirstProject/pkg/email"
	"MyFirstProject/pkg/utils/ctl"
	"MyFirstProject/pkg/utils/jwt"
	"MyFirstProject/pkg/utils/log"
	util "MyFirstProject/pkg/utils/upload"
	"MyFirstProject/repository/db/dao"
	"MyFirstProject/repository/db/model"
	"MyFirstProject/types"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// 用户注册

func (s *UserSrv) UserRegister(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if exist {
		err = errors.New("用户已经存在了")
		return
	}
	user := &model.User{
		NickName: req.NickName,
		UserName: req.UserName,
		Status:   model.Active,
		Money:    consts.UserInitMoney,
	}
	//加密密码
	if err = user.SetPassword(req.Password); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	//加密money
	money, err := user.EncryptMoney(req.Key)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user.Money = money
	// 默认头像走的local
	user.Avatar = consts.UserDefaultAvatarLocal
	if conf.Config.System.UploadModel == consts.UploadModelOss {
		// 如果配置是走oss，则用url作为默认头像
		user.Avatar = consts.UserDefaultAvatarOss
	}

	// 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	return
}

// 用户登录

func (s *UserSrv) UserLogin(ctx context.Context, req *types.UserServiceReq) (resp interface{}, err error) {
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if !exist {
		log.LogrusObj.Error(err)
		return nil, errors.New("用户不存在")
	}
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("账号/密码不正确")
	}
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	userResp := &types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   user.AvatarURL(),
		CreateAt: user.CreatedAt.Unix(),
	}

	resp = &types.UserTokenData{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}

//用户修改信息

func (s *UserSrv) UserInfoUpdate(ctx context.Context, req *types.UserInfoUpdateReq) (resp interface{}, err error) {
	//找到用户
	u, _ := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	err = userDao.UpdateUserById(u.Id, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}

//用户修改头像

func (s *UserSrv) UserAvatarUpload(ctx context.Context, file multipart.File, fileSize int64, req *types.UserServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	uId := u.Id
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal { // 兼容两种存储方式
		path, err = util.AvatarUploadToLocalStatic(file, uId, user.UserName)
	} else {
		path, err = util.UploadToQiNiu(file, fileSize)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	return
}

// 用户信息展示

func (s *UserSrv) UserInfoShow(ctx context.Context, req *types.UserInfoShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user, err := dao.NewUserDao(ctx).GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		NickName: user.NickName,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   user.AvatarURL(),
		CreateAt: user.CreatedAt.Unix(),
	}

	return
}

//发送邮件

func (s *UserSrv) SendEmail(ctx context.Context, req *types.SendEmailServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var address string
	token, err := jwt.GenerateEmailToken(u.Id, req.OperationType, req.Email, req.Password)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	sender := email.NewEmailSender()
	address = conf.Config.Email.ValidEmail + token
	mailText := fmt.Sprintf(consts.EmailOperationMap[req.OperationType], address)
	if err = sender.Send(mailText, req.Email, "FanOneMall"); err != nil {
		log.LogrusObj.Error(err)
		return
	}

	return
}
