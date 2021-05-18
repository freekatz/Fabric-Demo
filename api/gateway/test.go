
package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
)

type MyMux struct {
}

func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		sayhelloName(w, r)
	case "/login":
		login(w, r)
	default:
		http.NotFound(w, r)
	}
    return
}


func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()       // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
    // 注意:如果没有调用 ParseForm 方法，下面无法获取表单的数据
    fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) // 获取请求的方法
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        log.Println(t.Execute(w, nil))
    } else {
        // err := r.ParseForm()   // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
        // if err != nil {
        //    // handle error http.Error() for example
        //   log.Fatal("ParseForm: ", err)
        // }
        // 请求的是登录数据，那么执行登录的逻辑判断
        // fmt.Println("username:", r.Form["username"])
        // fmt.Println("password:", r.Form["password"])
		
		fmt.Println("username:", template.HTMLEscapeString(r.FormValue("username"))) // 输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.FormValue("password")))
		template.HTMLEscape(w, []byte(r.FormValue("username") + "\n")) // 输出到客户端
    }
}

func main() {
	// mux := &MyMux{}
    http.HandleFunc("/", sayhelloName)       // 设置访问的路由
    http.HandleFunc("/login", login)         // 设置访问的路由
    // err := http.ListenAndServe(":9090", mux) // 设置监听的端口
    err := http.ListenAndServe(":9090", nil) // 设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}