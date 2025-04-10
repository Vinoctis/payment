package middleware

import (
	"github.com/gin-gonic/gin"
	"payment/pkg/utils"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context){
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			utils.ResponseError(c, utils.CodeNeedLogin)
			c.Abort()
			return
		}
		c.Set("userID", 1)
		c.Next()
	}
}