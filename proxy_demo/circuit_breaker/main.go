package main

import (
	"github.com/donscoco/gateway/middleware_http"
	"github.com/donscoco/gateway/police"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	"log"
	"net/http"
)

var addr = ":8001"

// 熔断方案
func main() {
	lb := load_balance.LoadBalanceFactory(load_balance.LbRoundRobin)
	lb.Add("http://192.168.2.193:8002")
	lb.Add("http://192.168.2.193:8003")

	coreFunc := func(c *middleware_http.SliceRouterContext) http.Handler {
		return proxy.NewLoadBalanceReverseProxy(c, lb)
	}

	log.Println("Starting httpserver at " + addr)

	police.ConfCricuitBreaker(true) // 开启
	sliceRouter := middleware_http.NewSliceRouter()
	sliceRouter.Group("/").Use(middleware_http.CircuitMW()) // 配置中间件进行统计和升降级
	routerHandler := middleware_http.NewSliceRouterHandler(coreFunc, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}
