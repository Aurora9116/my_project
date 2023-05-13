package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/account"
	"time"
)

type HandlerAccount struct {
}

func (a HandlerAccount) account(c *gin.Context) {
	// 接收请求参数
	result := &common.Result{}
	var req *model.AccountReq
	_ = c.ShouldBind(&req)
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 调用grpc模块
	msg := &account.AccountReqMessage{
		MemberId:         memberId,
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		OrganizationCode: c.GetString("organizationCode"),
		DepartmentCode:   req.DepartmentCode,
		SearchType:       int32(req.SearchType),
	}
	response, err := AccountService.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	// 返回数据
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	}))
}

func NewAccount() *HandlerAccount {
	return &HandlerAccount{}
}
