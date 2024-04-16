package routes

import (
	api "MyFirstProject/api/v1"
	"MyFirstProject/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-sercret"))
	r.Use(middleware.Cors(), middleware.Jaeger())
	r.Use(sessions.Sessions("mysession", store))
	r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		// 用户操作
		//1.用户注册
		v1.POST("user/register", api.UserRegisterHandler())
		//2.用户登录
		v1.POST("user/login", api.UserLoginHandler())

		authed := v1.Group("/") // 需要登陆保护
		authed.Use(middleware.AuthMiddleware())
		{
			//用户信息修改
			authed.POST("user/update", api.UserUpdateHandler())
			//展示用户信息
			authed.GET("user/show_info", api.ShowUserInfoHandler())
			//用户修改头像
			authed.POST("user/avatar", api.UploadAvatarHandler())
			//发送验证码
			authed.POST("user/send_email", api.SendEmailHandler())
			//取消关注
			authed.POST("user/unfollowing", api.UserUnFollowingHandler())
			//关注
			authed.POST("user/following", api.UserFollowingHandler())
		}
	}
	return r
}
