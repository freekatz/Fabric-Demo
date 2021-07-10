/*
将链码封装为 http 接口
*/

package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/1uvu/fabric-sdk-client/pkg/client"
	"github.com/1uvu/fabric-sdk-client/pkg/types"
	"github.com/1uvu/serve"
)

//
// 实现 API Server, 定义路由, 处理逻辑等
//

// todo [暂不做] 将每一个路由的 handler 保存一一对应起来

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
	channelids map[string][]string = map[string][]string{
		"org1": []string{"channel1", "channel12", "channel123"},
		"org2": []string{"channel2", "channel12", "channel123"},
		"org3": []string{"channel3", "channel123"},
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

type Server struct{}

func (s *Server) Run(addr string) {
	log.Println("http api.")

	// 先初始化允许的操作, 否则下面调用接口无法提供任何功能
	initOPTs()

	// 删除 app client wallets
	removeAllWallet()

	g := serve.Default()

	{
		g.GET("/", func(c *serve.Context) {
			c.String(http.StatusOK, "api server start success.")
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	// app 接口可以操作 2 种客户端，只包含链码调用相关的操作
	app := g.Group("/app")
	app.Use(onlyForApp()) // app group middleware
	{
		app.POST("/:orgid/:channelid/chaincode/invoke", func(c *serve.Context) {
			data, _ := ioutil.ReadAll(c.Req.Body)
			request := new(types.InvokeRequest)
			json.Unmarshal(data, request)
			appInvokeHandle(c, request)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	// admin 接口可以操作 5 种客户端，提供各种操作
	admin := g.Group("/admin")
	admin.Use(onlyForAdmin()) // v1 group middleware
	{
		admin.POST("/:orgid/:channelid/:client/:opt", func(c *serve.Context) {
			data, _ := ioutil.ReadAll(c.Req.Body)
			request := new(types.InvokeRequest)
			json.Unmarshal(data, request)
			adminInvokeHandle(c, request)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	g.Run(addr)
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

func appInvokeHandle(c *serve.Context, request *types.InvokeRequest) {

	// 1 todo [暂时不做] 判断 url 的动态路由合法性: 是否在允许的范围内
	//=// c.Param("orgid")
	//=// c.Param("client")
	//=// c.Param("opt")

	for i := range request.Args {
		request.Args[i] = strings.Replace(request.Args[i], `'`, `"`, -1)
	}

	app, err := client.GetApp(c.Param("channelid"), configs[c.Param("orgid")])

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invoke failed with error: %s", err))
	}

	// 2 执行链码调用, 获取结果
	log.Println(request)
	resp, err := app.InvokeChaincode(request)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invoke failed with error: %s", err))
	}

	// 3 使用 c.JSON() 返回调用结果

	c.JSON(http.StatusOK, resp) // replace with a other handler

}

func adminInvokeHandle(c *serve.Context, request *types.InvokeRequest) {

	// 1 todo [暂时不做] 判断 url 的动态路由合法性: 是否在允许的范围内
	//=// c.Param("orgid")
	//=// c.Param("client")
	//=// c.Param("opt")

	for i := range request.Args {
		request.Args[i] = strings.Replace(request.Args[i], `'`, `"`, -1)
	}

	admin, err := client.GetAdmin(configs[c.Param("orgid")])

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invoke failed with error: %s", err))
	}

	// 2 执行链码调用, 获取结果

	app, err := admin.GetAppClient(c.Param("channelid"))

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invoke failed with error: %s", err))
	}

	log.Println(request)
	resp, err := app.InvokeChaincode(request)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invoke failed with error: %s", err))
	}

	// 3 使用 c.JSON() 返回调用结果

	c.JSON(http.StatusOK, resp) // replace with a other handler

}

//
// middlewares
//

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

//
// utils
//

func removeAllWallet() {
	_ = os.Remove("./wallet")
	_ = os.Remove("./keystore")
}
