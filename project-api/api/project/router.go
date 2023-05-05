package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/midd"
	"test.com/project-api/router"
)

func init() {
	log.Printf("init project router \n")
	ru := &RouterProject{}
	router.Register(ru)
}

type RouterProject struct {
}

func (*RouterProject) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()

	h := New()
	group := r.Group("/project")
	// todo bug
	group.Use(midd.TokenVerify())
	group.POST("/index", h.index)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)
	group.POST("/project/recovery", h.recoveryProject)
	group.POST("/project_collect/collect", h.collectProject)
	group.POST("/project/edit", h.editProject)

	t := NewTask()
	group.POST("/task_stages", t.taskStages)
}
