package main

import (
	"fmt"
	"github.com/donscoco/gateway/base_server/jwt"
	"github.com/donscoco/gateway/middleware_http"
	"github.com/donscoco/gateway/police/load_balance"
	"github.com/donscoco/gateway/proxy"
	"net/http"
	"strings"
)

func main() {
	// 1.

	//lb := load_balance.LoadBalanceFactory(load_balance.LbRoundRobin)
	//lb.Add("192.168.2.193:8002")
	lbconf, _ := load_balance.NewLoadBalanceZkConf(
		"http://%s",
		"/real_server",
		[]string{"192.168.2.132:2181"},
		map[string]string{"testkey": "testval"})
	lb := load_balance.LoadBalanceFactorWithConf(load_balance.LbRoundRobin, lbconf)

	router := middleware_http.NewSliceRouter()
	checkFunc := func(ctx *middleware_http.SliceRouterContext) bool {
		token := ctx.Req.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", -1)
		fmt.Println(token)
		msg, err := jwt.Decode(token)
		if err != nil {
			return false
		}
		// todo 校验逻辑
		fmt.Println("check valid:", msg)
		return true
	}
	router.Group("/").Use(middleware_http.JwtMiddleWare(checkFunc))
	router.Group("/get_token").Use(func(c *middleware_http.SliceRouterContext) {
		user := ""
		pwd := ""
		c.Req.ParseForm()
		if len(c.Req.Form["user"]) > 0 {
			user = c.Req.Form["user"][0]
		}
		if len(c.Req.Form["pwd"]) > 0 {
			pwd = c.Req.Form["pwd"][0]
		}
		fmt.Println("user", user)
		fmt.Println("pwd", pwd)
		if user == "ironhead" && pwd == "123456" {
			jwtToken, err := jwt.Encode(user)
			if err != nil {
				c.Rw.Write([]byte("get token error:" + err.Error()))
				c.Abort()
				return
			}
			c.Rw.Write([]byte(jwtToken))
			c.Abort()
			return
		}
		c.Rw.Write([]byte("get token error:wrong user or secret"))
		c.Abort()
		return
	})
	coreFunc := func(c *middleware_http.SliceRouterContext) http.Handler {
		return proxy.NewLoadBalanceReverseProxy(c, lb)
	}
	handler := middleware_http.NewSliceRouterHandler(coreFunc, router)

	http.ListenAndServe(":8001", handler)

}
