package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"
	_ "test.com/project-api/api"
	"test.com/project-api/api/midd"
	"test.com/project-api/config"
	"test.com/project-api/router"
	"test.com/project-api/tracing"
	srv "test.com/project-common"
)

func main() {
	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r.Use(midd.RequestLog())
	r.Use(otelgin.Middleware("project-api"))
	r.StaticFS("/upload", http.Dir("upload"))
	router.InitRouter(r)
	// 开启pprof 默认访问路径是/debug/pprof 或自定义 例如 pprof.Register(r, "/info/pprof") 路径即为 /info/pprof
	pprof.Register(r)
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}
