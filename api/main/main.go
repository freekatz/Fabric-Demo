/*
将链码封装为 http 接口
*/

package main

import (
	"log"
	"net/http"

	"github.com/1uvu/Fabric-Demo/api/gateway"
)

func main() {
	log.Println("http api.")

	g := gateway.New()
	g.GET("/", func(c *gateway.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	g.GET("/hello", func(c *gateway.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	g.POST("/login", func(c *gateway.Context) {
		c.JSON(http.StatusOK, gateway.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	g.Run(":9999")
}
