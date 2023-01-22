package main

import (
	//"github.com/donscoco/gateway/base_server"
	"github.com/donscoco/gateway/middleware_tcp"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	server_tcp "github.com/donscoco/gateway/server/tcp"
	"log"
)

func main() {

	//1. 创建中间件用到的基础服务

	//2. 创建负载均衡策略

	//3. 创建路由 并 创建handler 并创建 代理

	//4. 启动tcp server 服务器

	//localcounter, _ := base_server.NewFlowCountService("my_app", 10*time.Second)

	//lbconf, _ := load_balance.NewLoadBalanceZkConf(
	//	"%s", "/real_server",
	//	[]string{"192.168.2.132:2181"},
	//	map[string]string{"": ""},
	//)
	//lb := load_balance.LoadBalanceFactorWithConf(load_balance.LbRoundRobin, lbconf)

	lb := load_balance.LoadBalanceFactory(load_balance.LbRoundRobin)
	//lb.Add("127.0.0.1:6379") // 代理redis
	//lb.Add("127.0.0.1:30801") // 代理thrift
	//lb.Add("127.0.0.1:30701") // 代理grpc, grpc 代理不了，todo：做grpc 代理

	router := middleware_tcp.NewTcpSliceRouter()
	router.Group("/").Use(
		//middleware_tcp.FlowCountLocal(localcounter),
		middleware_tcp.TraceLogSliceMW(),
	)
	// 调用链的最终处理(在这里就是负载均衡的处理逻辑)
	endpointFunc := func(c *middleware_tcp.TcpSliceRouterContext) server_tcp.TCPHandler {
		return proxy.NewTcpLoadBalanceReverseProxy(c, lb)
	}
	tcpHandler := middleware_tcp.NewTcpSliceRouterHandler(endpointFunc, router)

	log.Fatal(server_tcp.ListenAndServe(":7001", tcpHandler))
}
