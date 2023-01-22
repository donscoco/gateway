package middleware_tcp

import (
	"fmt"
	"github.com/donscoco/gateway/base_server"
)

func FlowCountLocal(counter *base_server.FlowCountService) func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		counter.Increase()
		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)
		c.Next()
	}
}
