package main

import (
	"github.com/donscoco/gateway/base_server"
	"github.com/donscoco/gateway/middleware_http"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	"net/http"
	"time"
)

func main() {

	// 1.先创建要用到的基础组件 baseserever: localcounter

	// 2.创建负载均衡策略 police:loadbalance

	// 3.创建中间件router middleware: logger,counter,
	// 3.1 创建代理 proxy:http 并 封装代理为 中间件router 的 handler
	// 3.2 创建中间件 handler

	// 5.server: http

	localCounter, _ := base_server.NewFlowCountService("my-http-proxy", 10*time.Second)

	loadbalanceConf, _ := load_balance.NewLoadBalanceZkConf("http://%s",
		"/real_server",
		[]string{"192.168.2.132:2181"},
		map[string]string{"127.0.0.1:2003": "20"})
	//lb := load_balance.LoadBanlanceFactory(load_balance.LbRoundRobin)
	lb := load_balance.LoadBalanceFactorWithConf(load_balance.LbRoundRobin, loadbalanceConf)

	router := middleware_http.NewSliceRouter()
	router.Group("/").Use(
		middleware_http.TraceLogSliceMW(),
		middleware_http.RateLimiter(),
		middleware_http.FlowCountLocal(localCounter),
	)
	coreFunc := func(c *middleware_http.SliceRouterContext) http.Handler {
		return proxy.NewLoadBalanceReverseProxy(c, lb)
	}
	routerHandler := middleware_http.NewSliceRouterHandler(coreFunc, router)

	err := http.ListenAndServe(":8001", routerHandler)
	if err != nil {
		panic(err)
	}
}
