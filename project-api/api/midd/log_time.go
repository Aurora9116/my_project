package midd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func RequestLog() func(*gin.Context) {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		diff := time.Now().UnixMilli() - start.UnixMilli()
		diff1 := time.Now().Unix() - start.Unix()
		diff2 := time.Now().UnixMicro() - start.UnixMicro()
		zap.L().Info(fmt.Sprintf("%s 用时 %d ms", c.Request.RequestURI, diff))
		zap.L().Info(fmt.Sprintf("%s 用时 %d ms", c.Request.RequestURI, diff1))
		zap.L().Info(fmt.Sprintf("%s 用时 %d ms", c.Request.RequestURI, diff2))
	}
}
