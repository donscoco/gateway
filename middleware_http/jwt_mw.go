package middleware_http

func JwtMiddleWare(handleFunc func(c *SliceRouterContext) bool) func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {

		if handleFunc(c) {
			c.Next()
		} else {
			c.Rw.Write([]byte("jwt auth invalid:"))
			c.Abort()
			return
		}

		//token := c.Req.Header.Get("Authorization")
		//token = strings.Replace(token, "Bearer ", "", -1)
		//if _, err := base_server.Decode(token); err != nil {
		//	c.Rw.Write([]byte("jwt auth invalid:" + err.Error()))
		//	c.Abort()
		//	return
		//}
		//c.Next()
	}
}
