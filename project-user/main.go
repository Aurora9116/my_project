package main

import (
	"github.com/gin-gonic/gin"
	_ "test.com/project-api/api"
	srv "test.com/project-common"
	"test.com/project-user/config"
	"test.com/project-user/router"
)

func main() {
	r := gin.Default()

	router.InitRouter(r)

	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	var stop = func() {
		gc.Stop()
	}

	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
