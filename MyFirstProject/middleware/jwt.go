package middleware

import (
	"MyFirstProject/consts"
	"MyFirstProject/pkg/e"
	"MyFirstProject/pkg/utils/ctl"
	util "MyFirstProject/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware token验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		var code int
		code = e.SUCCESS
		accessToken := context.GetHeader("access_token")
		refreshToken := context.GetHeader("refresh_token")
		if accessToken == "" {
			code = e.InvalidParams
			context.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Token不能为空",
			})
			context.Abort()
			return
		}
		newAccessToken, newRefreshToken, err := util.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		}
		if code != e.SUCCESS {
			context.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			context.Abort()
			return
		}
		claims, err := util.ParseToken(accessToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			context.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   err.Error(),
			})
			context.Abort()
			return
		}
		println(newRefreshToken)
		println(claims.ID)
		println(claims.Username)
		SetToken(context, newAccessToken, newRefreshToken)
		context.Request = context.Request.WithContext(ctl.NewContext(context.Request.Context(), &ctl.UserInfo{Id: claims.ID}))
		ctl.InitUserInfo(context.Request.Context())
		context.Next()
	}
}
func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
