package midd

import (
	"context"
	"github.com/gin-gonic/gin"
	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
	"time"
)

// GetIp 获取ip函数
func GetIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}
func TokenVerify() func(c *gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		// 1.从header中获取token
		token := c.GetHeader("Authorization")
		// 2.调用user服务进行token认证
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		ip := GetIp(c)
		response, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token, Ip: ip})

		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(200, result.Fail(code, msg))
			c.Abort()
			return
		}
		// 3.处理结束，认证通过 将信息放入gin的上下文 失败返回未登录
		c.Set("memberId", response.Member.Id)
		c.Set("memberName", response.Member.Name)
		c.Set("organizationCode", response.Member.OrganizationCode)
		c.Next()
	}
}
