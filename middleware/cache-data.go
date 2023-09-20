package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"one-api/common"
)

func CacheDataMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rdb := common.RDB
		// 尝试从缓存中获取数据
		ctx := context.Background()
		cachedData, err := rdb.Get(ctx, key).Result()
		if err == nil {
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cachedData))
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"data": "没有找到新数据"})
			return
		}
		c.Next()
	}
}
