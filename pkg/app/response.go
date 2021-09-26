package app

import (
	"onlineShopping/pkg/e"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, errCode int, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
	return
}
