package serve

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// current time
		now := time.Now()
		// process request
		c.Next()
		// logging
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(now))
	}
}
