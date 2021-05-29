/*
将链码封装为 http 接口
*/

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/1uvu/Fabric-Demo/api/gateway"
)

func onlyForV2() gateway.HandlerFunc {
	// 执行顺序: onlyForV2 -> handle -> logger
	return func(c *gateway.Context) {
		// Start timer
		t := time.Now()
		// // if a server error occurred
		// c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	log.Println("http api.")

	g := gateway.New()
	g.Use(gateway.Logger()) // global midlleware
	g.GET("/", func(c *gateway.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gateway</h1>")
	})

	g.GET("/index", func(c *gateway.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
		log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
	})
	v1 := g.Group("/v1")
	{
		v1.GET("/", func(c *gateway.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gateway</h1>")
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})

		v1.GET("/hello", func(c *gateway.Context) {
			// expect /hello?name=zjh
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}
	v2 := g.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gateway.Context) {
			// expect /hello/zjh
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
		v2.POST("/login", func(c *gateway.Context) {
			c.JSON(http.StatusOK, gateway.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})

	}

	g.Run(":9999")
}
