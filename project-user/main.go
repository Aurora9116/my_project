package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	_ "test.com/project-api/api"
	srv "test.com/project-common"
	"test.com/project-user/config"
	"test.com/project-user/router"
	"test.com/project-user/tracing"
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

	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	var stop = func() {
		gc.Stop()
	}

	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
