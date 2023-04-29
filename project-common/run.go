package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		log.Printf("%s running in %s \n", srvName, server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR字符(Ctrl + C) 触发
	// SIGTERM结束程序（可以被捕获、阻塞或忽略）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutting Down project %s ...\n", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v\n", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Printf("wait timeout...\n")
	}
	log.Printf("%s stop success...\n", srvName)
}
