package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	_ "test.com/project-api/api"
	srv "test.com/project-common"
	"test.com/project-project/config"
	"test.com/project-project/router"
	"test.com/project-project/tracing"
)

func main() {
	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	router.InitRouter(r)

	// 初始化rpc调用
	router.InitUserRpc()
	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	// 初始化kafka
	Close := config.InitKafkaWriter()
	var stop = func() {
		gc.Stop()
		Close()
	}
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
