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
	userRoute := apiRouter.Group("/v1")
	{
		//只有chat需要限制
		userRoute.POST("/chat", middleware.GlobalAPIRateLimit(), controller.RelayChatbase)
	}
	emailRouter := router.Group("/email")
	{
		emailRouter.POST("/send", controller.SendEmail)
	}
}
