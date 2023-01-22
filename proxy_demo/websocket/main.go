package main

import (
	"context"
	"github.com/donscoco/gateway/middleware_http"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	"log"
	"net/http"
)

var (
	addr = "127.0.0.1:20001"
)

func main() {
	// todo 改成 zk 注册服务发现
	rb := load_balance.LoadBalanceFactory(load_balance.LbWeightRoundRobin)
	rb.Add("http://127.0.0.1:20002", "50")
	rb.Add("http://127.0.0.1:20003", "50")

	proxy := proxy.NewLoadBalanceReverseProxy(&middleware_http.SliceRouterContext{Ctx: context.Background()}, rb)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
