package v1

import (
	"MyFirstProject/consts"
	"MyFirstProject/pkg/utils/ctl"
	"MyFirstProject/pkg/utils/log"
	"MyFirstProject/service"
	"MyFirstProject/types"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

//用户注册

func UserRegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		//参数校验
		if req.Key == "" || len(req.Key) != consts.EncryptMoneyKeyLength {
			err := errors.New("key长度错误,必须是6位数")
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserRegister(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// 用户登录
func UserLoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserServiceReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(ctx, err))
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserLogin(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
