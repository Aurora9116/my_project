package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

//
//func RegisterGrpc() *grpc.Server {
//	c := gRPCConfig{
//		Addr: config.C.GC.Addr,
//		RegisterFunc: func(g *grpc.Server) {
//			login.RegisterLoginServiceServer(g, loginServiceV1.New())
//		},
//	}
//	s := grpc.NewServer()
//	c.RegisterFunc(s)
//	lis, err := net.Listen("tcp", c.Addr)
//	if err != nil {
//		log.Println("cannot listen")
//	}
//	go func() {
//		err = s.Serve(lis)
//		if err != nil {
//			log.Println("server started error", err)
//			return
//		}
//	}()
//	return s
//}
