package middleware_tcp

import "log"

func TraceLogSliceMW() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		log.Println("trace_in")
		c.Next()
		log.Println("trace_out")
	}
}
