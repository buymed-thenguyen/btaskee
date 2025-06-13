package middleware

import (
	"btaskee/model/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() {
			return
		}

		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			Error(c, status, c.Errors.String())
		} else {
			Success(c)
		}
	}
}

func Success(c *gin.Context) {
	data, _ := c.Get(constant.DATA_CTX)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}
