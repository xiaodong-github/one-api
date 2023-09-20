package router

import (
	"github.com/gin-gonic/gin"
	"one-api/controller"
	"one-api/middleware"
)

func SetChatbaseRouter(router *gin.Engine) {
	// 聊天机器人
	relayV1Router := router.Group("/cb")
	relayV1Router.Use(middleware.TokenAuth(), middleware.Distribute())
	{
		relayV1Router.POST("/completions", controller.Relay)
		relayV1Router.POST("/chat/completions", controller.Relay)
		relayV1Router.POST("/chat", controller.RelayChatbase)
		relayV1Router.POST("/edits", controller.Relay)
		relayV1Router.POST("/images/generations", controller.Relay)
		relayV1Router.POST("/images/edits", controller.Relay)
		relayV1Router.POST("/images/variations", controller.RelayNotImplemented)
		relayV1Router.POST("/embeddings", controller.Relay)
		relayV1Router.POST("/engines/:model/embeddings", controller.Relay)
		relayV1Router.POST("/audio/transcriptions", controller.RelayNotImplemented)
		relayV1Router.POST("/audio/translations", controller.RelayNotImplemented)
		relayV1Router.GET("/files", controller.RelayNotImplemented)
		relayV1Router.POST("/files", controller.RelayNotImplemented)
		relayV1Router.DELETE("/files/:id", controller.RelayNotImplemented)
		relayV1Router.GET("/files/:id", controller.RelayNotImplemented)
		relayV1Router.GET("/files/:id/content", controller.RelayNotImplemented)
		relayV1Router.POST("/fine-tunes", controller.RelayNotImplemented)
		relayV1Router.GET("/fine-tunes", controller.RelayNotImplemented)
		relayV1Router.GET("/fine-tunes/:id", controller.RelayNotImplemented)
		relayV1Router.POST("/fine-tunes/:id/cancel", controller.RelayNotImplemented)
		relayV1Router.GET("/fine-tunes/:id/events", controller.RelayNotImplemented)
		relayV1Router.DELETE("/models/:model", controller.RelayNotImplemented)
		relayV1Router.POST("/moderations", controller.Relay)
	}
}
