package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "test.com/project-api/api"
	"test.com/project-api/api/midd"
	"test.com/project-api/config"
	"test.com/project-api/router"
	srv "test.com/project-common"
)

func main() {
	r := gin.Default()
	r.Use(midd.RequestLog())
	r.StaticFS("/upload", http.Dir("upload"))
	router.InitRouter(r)
	// 开启pprof 默认访问路径是/debug/pprof 或自定义 例如 pprof.Register(r, "/info/pprof") 路径即为 /info/pprof
	pprof.Register(r)
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}
