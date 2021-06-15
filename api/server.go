/*
将链码封装为 http 接口
*/

package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/1uvu/fabric-sdk-client/types"
	"github.com/1uvu/serve"
)

// todo 将每一个路由的 handler 保存一一对应起来
var (
	orgids []string = []string{
		"org1",
		"org2",
		"org3",
	}
	configs map[string]string = map[string]string{
		"org1": "./config/client-org1.yaml",
		"org2": "./config/client-org2.yaml",
		"org3": "./config/client-org3.yaml",
	}
	clients []string = []string{
		"chaincode", // 代替 channel 和 event client, 主要是 chaincode 相关的操作
		"ledger",    // 主要是 ledger 的原始访问操作
		"resource",  // 主要是 channel 相关的管理操作
		"msp",       // 主要是 msp 成员关系管理的操作
	}
	// 每种 client 对应的 opts
	opts map[string][]string = make(map[string][]string)
)

func Run() {
	log.Println("http api.")

	// 先初始化允许的操作, 否则下面调用接口无法提供任何功能
	initOPTs()

	g := serve.Default()
	g.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	// app 接口可以操作 2 种客户端，只包含链码调用相关的操作
	app := g.Group("/app")
	app.Use(onlyForApp()) // app group middleware
	{
		app.POST("/app/:orgid/chaincode/invoke", func(c *serve.Context) {
			// params := &types.InvokeParams{
			params := &types.InvokeParams{
				ChaincodeID: c.PostForm("chaincodeID"),
				Fcn:         c.PostForm("fcn"),
				Args:        str2strArray(c.PostForm("args")),
				NeedSubmit:  str2bool(c.PostForm("needSubmit")),
			}
			appInvokeHandle(c, params)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	// admin 接口可以操作 5 种客户端，提供各种操作
	admin := g.Group("/admin")
	admin.Use(onlyForAdmin()) // v1 group middleware
	{
		admin.POST("/admin/:orgid/:client/:opt", func(c *serve.Context) {
			// params := &types.InvokeParams{
			params := &types.InvokeParams{
				ChaincodeID: c.PostForm("chaincodeID"),
				Fcn:         c.PostForm("fcn"),
				Args:        str2strArray(c.PostForm("args")),
				NeedSubmit:  str2bool(c.PostForm("needSubmit")),
				Endpoints:   str2strArray(c.PostForm("endpoints")),
			}
			adminInvokeHandle(c, params)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	g.Run(":9999")
}

// define client's opts
func initOPTs() {
	opts["chaincode"] = []string{
		"invoke",
		// add more opts for chaincode client
	}

	// only for admin
	opts["ledger"] = []string{
		"queryChannelInfo",
		"queryBlockInfo",
		// add more opts for chaincode client
	}

	// define other client's opts
}

func onlyForApp() serve.HandlerFunc {
	// 执行顺序: onlyForApp -> handle -> logger
	return func(c *serve.Context) {
		// Start timer
		t := time.Now()
		// // if a server error occurred
		// c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group app", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func onlyForAdmin() serve.HandlerFunc {
	// 执行顺序: onlyForAdmin -> handle -> logger
	return func(c *serve.Context) {
		// Start timer
		t := time.Now()
		// // if a server error occurred
		// c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group admin", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

// func appInvokeHandle(c *serve.Context, params *types.InvokeParams) {
func appInvokeHandle(c *serve.Context, params *types.InvokeParams) {

	c.JSON(http.StatusOK, params) // replace with a other handler

	// 1 todo [暂时不做] 判断 url 的动态路由合法性: 是否在允许的范围内
	//=// c.Param("orgid")
	//=// c.Param("client")
	//=// c.Param("opt")

	// 2 执行链码调用, 获取结果

	// 3 使用 c.JSON() 返回调用结果

}

// func adminInvokeHandle(c *serve.Context, params *types.InvokeParams) {
func adminInvokeHandle(c *serve.Context, params *types.InvokeParams) {

	c.JSON(http.StatusOK, params) // replace with a other handler

	// 1 todo [暂时不做] 判断 url 的动态路由合法性: 是否在允许的范围内
	//=// c.Param("orgid")
	//=// c.Param("client")
	//=// c.Param("opt")

	// 2 执行链码调用, 获取结果

	// 3 使用 c.JSON() 返回调用结果

}

func str2strArray(str string) []string {
	return []string{}
}

func str2bool(str string) bool {
	return false
}

// type person struct {
// 	Name string
// 	Age  int8
// }

// func onlyForV2() serve.HandlerFunc {
// 	// 执行顺序: onlyForV2 -> handle -> logger
// 	return func(c *serve.Context) {
// 		// Start timer
// 		t := time.Now()
// 		// // if a server error occurred
// 		// c.Fail(500, "Internal Server Error")
// 		// Calculate resolution time
// 		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
// 	}
// }

// func FormatAsDate(t time.Time) string {
// 	year, month, day := t.Date()
// 	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
// }

// func main() {
// 	log.Println("http api.")

// 	g := serve.Default()
// 	g.SetFuncMap(template.FuncMap{
// 		"FormatAsDate": FormatAsDate,
// 	})
// 	g.LoadHTMLGlob("templates/*")
// 	g.Static("/assets", "./assets")
// 	// g.Static("/assets/css", "./css") // 会覆盖之前的路由中 `重复` 的文件，所有为了避免出现异常，建议谨慎设计路由

// 	p1 := &person{Name: "zjh", Age: 22}
// 	p2 := &person{Name: "qx", Age: 20}

// 	g.GET("/", func(c *serve.Context) {
// 		c.HTML(http.StatusOK, "css.html", serve.H{
// 			"title": "CSS",
// 		})
// 	})

// 	g.GET("/index", func(c *serve.Context) {
// 		c.HTML(http.StatusOK, "<h1>Index Page</h1>", nil)
// 		log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
// 	})

// 	g.GET("/panic", func(c *serve.Context) {
// 		texts := []string{"panic"}
// 		c.String(http.StatusOK, texts[1])
// 	})

// 	g.GET("/person", func(c *serve.Context) {
// 		c.HTML(http.StatusOK, "array.html", serve.H{
// 			"title":       "Persons",
// 			"personArray": [2]*person{p1, p2},
// 		})
// 	})

// 	g.GET("/date", func(c *serve.Context) {
// 		c.HTML(http.StatusOK, "func.html", serve.H{
// 			"title": "Function Date",
// 			"now":   time.Date(2021, 5, 30, 0, 0, 0, 0, time.UTC),
// 		})
// 	})

// 	v1 := g.Group("/v1")
// 	{
// 		v1.GET("/", func(c *serve.Context) {
// 			c.HTML(http.StatusOK, "<h1>Hello serve</h1>", nil)
// 			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
// 		})

// 		v1.GET("/hello", func(c *serve.Context) {
// 			// expect /hello?name=zjh
// 			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
// 			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
// 		})
// 	}
// 	v2 := g.Group("/v2")
// 	v2.Use(onlyForV2()) // v2 group middleware
// 	{
// 		v2.GET("/hello/:name", func(c *serve.Context) {
// 			// expect /hello/zjh
// 			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
// 			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
// 		})
// 		v2.POST("/login", func(c *serve.Context) {
// 			c.JSON(http.StatusOK, serve.H{
// 				"username": c.PostForm("username"),
// 				"password": c.PostForm("password"),
// 			})
// 			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
// 		})

// 	}

// 	g.Run(":9999")
// }
