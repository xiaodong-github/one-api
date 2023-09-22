package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"one-api/controller"
	"one-api/middleware"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	//gzip 压缩中间件
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	//在 apiRouter 路由组中应用全局 API 速率限制中间件
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		userRoute := apiRouter.Group("/v1")
		{
			userRoute.POST("/chat", controller.RelayChatbase)
		}
	}
}
