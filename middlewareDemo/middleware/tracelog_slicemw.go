package middleware

import "log"

func TraceLogSliceMW() func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		log.Println("trace_in")
		defer log.Println("trace_on")
		c.Next()
	}
}
