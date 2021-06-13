/*
将链码封装为 http 接口
*/

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	// "github.com/1uvu/fabric-sdk-client/types"
	"github.com/1uvu/serve"
	// "github.com/1uvu/Fabric-Demo/api"
)

// temp
// chaincode invoke params
type InvokeParams struct {
	ChaincodeID string
	Fcn         string
	Args        []string
	NeedSubmit  bool
	// for admin client
	Endpoints []string
}

// todo 将每一个路由的 handler 保存一一对应起来
var (
	orgids []string = []string{
		"org1",
		"org2",
		"org3",
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

func initOPTs() {
	opts["chaincode"] = []string{
		"invoke",
		// add more opts for chaincode client
	}

	// define other client's opts
}

func onlyForV1() serve.HandlerFunc {
	// 执行顺序: onlyForV2 -> handle -> logger
	return func(c *serve.Context) {
		// Start timer
		t := time.Now()
		// // if a server error occurred
		// c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	log.Println("http api.")

	// 先初始化允许的操作, 否则下面调用接口无法提供任何功能
	initOPTs()

	g := serve.Default()
	g.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	v1 := g.Group("/v1")
	v1.Use(onlyForV1()) // v1 group middleware
	{
		v1.POST("/app/:orgid/:client/:opt", func(c *serve.Context) {
			// params := &types.InvokeParams{
			params := &InvokeParams{
				ChaincodeID: c.PostForm("chaincodeID"),
				Fcn:         c.PostForm("fcn"),
				Args:        str2strArray(c.PostForm("args")),
				NeedSubmit:  str2bool(c.PostForm("needSubmit")),
			}
			invokeHandle(c, params)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
		v1.POST("/admin/:orgid/:client/:opt", func(c *serve.Context) {
			// params := &types.InvokeParams{
			params := &InvokeParams{
				ChaincodeID: c.PostForm("chaincodeID"),
				Fcn:         c.PostForm("fcn"),
				Args:        str2strArray(c.PostForm("args")),
				NeedSubmit:  str2bool(c.PostForm("needSubmit")),
				Endpoints:   str2strArray(c.PostForm("endpoints")),
			}
			invokeHandle(c, params)
			log.Printf("=[Status Code: %d]=[Method: %4s]=[Path: %6s]\n", c.StatusCode, c.Method, c.Path)
		})
	}

	g.Run(":9999")
}

// func invokeHandle(c *serve.Context, params *types.InvokeParams) {
func invokeHandle(c *serve.Context, params *InvokeParams) {

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
