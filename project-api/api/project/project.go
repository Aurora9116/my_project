package project

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"test.com/project-api/pkg/model"
	"test.com/project-api/pkg/model/menu"
	"test.com/project-api/pkg/model/pro"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
	"time"
)

type HandlerProject struct {
}

func New() *HandlerProject {
	return &HandlerProject{}
}

func (p *HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}

	indexResponse, err := Pro.Index(ctx, msg)

	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	menus := indexResponse.Menus
	var ms []*menu.Menu
	copier.Copy(&ms, menus)
	c.JSON(http.StatusOK, result.Success(ms))
}

func (p *HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		Page:       page.Page,
		PageSize:   page.PageSize,
		SelectBy:   selectBy,
	}
	myProjectResponse, err := Pro.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": myProjectResponse.Total,
	}))
}

func (p *HandlerProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		Page:             page.Page,
		PageSize:         page.PageSize,
		ViewType:         int32(viewType),
		OrganizationCode: c.GetString("organizationCode"),
	}
	templateReponse, err := Pro.FindProjectTemplate(ctx, msg)
	fmt.Println("templateReponse==>", templateReponse)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var pms []*pro.ProjectTemplate
	copier.Copy(&pms, templateReponse.Ptm)
	if pms == nil {
		pms = []*pro.ProjectTemplate{}
	}
	for _, v := range pms {
		if v.TaskStages == nil {
			v.TaskStages = []*pro.TaskStagesOnlyName{}
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": templateReponse.Total,
	}))
}
func (p *HandlerProject) projectSave(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	var req *pro.SaveProjectRequest
	c.ShouldBind(&req)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		Name:             req.Name,
		TemplateCode:     req.TemplateCode,
		Description:      req.Description,
		Id:               int64(req.Id),
	}
	saveProject, err := Pro.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var rsp *pro.SaveProject
	copier.Copy(&rsp, saveProject)
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (p *HandlerProject) ReadProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detail, err := Pro.FindProjectDetail(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	pd := &pro.ProjectDetail{}
	copier.Copy(&pd, detail)
	c.JSON(http.StatusOK, result.Success(pd))
}
