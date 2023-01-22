package main

import (
	//"context"
	"fmt"
	"log"
	"net"
	//"time"

	//"github.com/donscoco/gateway/base_server"
	//"github.com/donscoco/gateway/middleware_grpc"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	"github.com/donscoco/gateway/proxy/grpc_proxy_ext"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/metadata"
	//"google.golang.org/grpc/status"
	//"strings"
)

const port = ":7001"

func main() {

	// 1.创建监听端口fd
	// 2.创建策略
	// 3.创建中间件要用到的基础组件
	// 4。

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rb := load_balance.LoadBalanceFactory(load_balance.LbWeightRoundRobin)
	rb.Add("localhost:30801", "40")

	//counter, _ := base_server.NewFlowCountService("local_app", time.Second)

	grpcHandler := proxy.NewGrpcLoadBalanceHandler(rb)
	s := grpc.NewServer(
		//grpc.ChainUnaryInterceptor( // 用于单元的
		//	middleware_grpc.GrpcAuthUnaryInterceptor,
		//	middleware_grpc.GrpcFlowCountUnaryInterceptor,
		//),
		//grpc.ChainStreamInterceptor( // 用于流式的
		//	middleware_grpc.GrpcAuthStreamInterceptor,
		//	middleware_grpc.GrpcFlowCountStreamInterceptor(counter),
		//),
		grpc.CustomCodec(grpc_proxy_ext.Codec()),
		//grpc.UnknownServiceHandler(grpc_proxy_ext.TransparentHandler(director)))
		//grpc.UnknownServiceHandler(routerHandler), // todo 可以不使用gpc自带的interceptor，自己实现一个链路中间件，类似tcp中间件那样做
		grpc.UnknownServiceHandler(grpcHandler),
	)

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
