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
		//商品操作
		//1.轮播图
		v1.GET("carousels", api.ListCarouselsHandler()) // 轮播图
		//2.查看商品
		v1.GET("product/list", api.ListProductsHandler())
		//3.获取商品详情
		v1.GET("product/show", api.ShowProductHandler())
		//4.搜索商品
		v1.POST("product/search", api.SearchProductsHandler())
		//5.商品图片
		v1.GET("product/imgs/list", api.ListProductImgHandler()) // 商品图片

		authed := v1.Group("/") // 需要登陆保护
		authed.Use(middleware.AuthMiddleware())
		{
			//用户操作
			//1.用户信息修改
			authed.POST("user/update", api.UserUpdateHandler())
			//2.展示用户信息
			authed.GET("user/show_info", api.ShowUserInfoHandler())
			//3.用户上传头像
			authed.POST("user/avatar", api.UploadAvatarHandler())
			//4.发送验证码
			authed.POST("user/send_email", api.SendEmailHandler())
			//5.取消关注
			authed.POST("user/unfollowing", api.UserUnFollowingHandler())
			//6.关注
			authed.POST("user/following", api.UserFollowingHandler())
			//7.邮箱验证
			authed.GET("user/valid_email", api.ValidEmailHandler())
			//收货地址操作
			//1.增加地址
			authed.POST("addresses/create", api.CreateAddressHandler())
			//2.展示某个地址
			authed.GET("addresses/show", api.ShowAddressHandler())
			//3.展示所有地址
			authed.GET("addresses/list", api.ListAddressHandler())
			//4.修改收货地址
			authed.POST("addresses/update", api.UpdateAddressHandler())
			//5.删除收货地址
			authed.POST("addresses/delete", api.DeleteAddressHandler())
			//购物车
			//1.增加购物车
			authed.POST("carts/create", api.CreateCartHandler())
			//2.查看购物车
			authed.GET("carts/list", api.ListCartHandler())
			//3.修改购物车
			authed.POST("carts/update", api.UpdateCartHandler()) // 购物车id
			//4.删除购物车
			authed.POST("carts/delete", api.DeleteCartHandler())
			//商品操作
			//1.增加商品
			authed.POST("product/create", api.CreateProductHandler())
			//2.删除商品
			authed.POST("product/delete", api.DeleteProductHandler())
			//3.更新商品
			authed.POST("product/update", api.UpdateProductHandler())
			//4.商品分类
			v1.GET("category/list", api.ListCategoryHandler())
			//收藏夹
			//1.创建收藏夹
			authed.POST("favorites/create", api.CreateFavoriteHandler())
			//2.获取收藏夹详情
			authed.GET("favorites/list", api.ListFavoritesHandler())
			//3.删除收藏夹
			authed.POST("favorites/delete", api.DeleteFavoriteHandler())
			//金额
			//1.显示金额
			authed.POST("money", api.ShowMoneyHandler())
			//订单操作
			//1.添加订单
			authed.POST("orders/create", api.CreateOrderHandler())
			//2.查看订单列表
			authed.GET("orders/list", api.ListOrdersHandler())
		}
	}
	return r
}
