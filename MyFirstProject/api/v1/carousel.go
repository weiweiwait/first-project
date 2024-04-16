package v1

import (
	"MyFirstProject/pkg/utils/ctl"
	"MyFirstProject/pkg/utils/log"
	"MyFirstProject/service"
	"MyFirstProject/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

//这段代码是一个用 Go 语言编写的中间件函数，适用于 Gin 框架，功能是处理获取轮播图列表的请求。
//以下是函数的详细步骤和功能：
//首先，该函数开始时创建了一个ListCarouselReq 类型的变量 req。ListCarouselReq 类型可能是一个用于存储请求中需要的数据的结构体。
//使用 ctx.ShouldBind() 函数尝试从请求中获取必要的参数并绑定到 req 变量。如果获取或绑定失败（如参数缺失或参数类型错误），则记录错误，向请求者返回错误信息，并结束函数处理。
//接着，通过 service.GetCarouselSrv() 函数获取轮播图服务对象 l。
//使用轮播图服务对象的 ListCarousel() 方法处理获取轮播图列表请求，并获取请求结果 resp 和可能的错误 err。如果处理请求时发生错误，同样记录错误信息，向请求者返回错误，并结束处理。
//最后，如果没有错误，就以 JSON 格式将成功的请求结果返回给请求者。
//整个函数处理流程是典型的 web 请求处理流程，接收请求，处理请求，返回结果。错误处理和日志记录也很到位。

func ListCarouselsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ListCarouselReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetCarouselSrv()
		resp, err := l.ListCarousel(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
