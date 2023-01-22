package middleware_http

import (
	"fmt"
	"github.com/donscoco/gateway/base_server"
)

func FlowCountLocal(counter *base_server.FlowCountService) func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		counter.Increase()
		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)
		c.Next()
	}
}
