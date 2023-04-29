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
	group := r.Group("/project/index")
	// todo bug
	group.Use(midd.TokenVerify())
	group.POST("", h.index)
	group1 := r.Group("/project/project")
	group1.Use(midd.TokenVerify())
	group1.POST("/selfList", h.myProjectList)
}
