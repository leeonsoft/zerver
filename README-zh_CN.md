[English](README.md)

#### Zerver-REST
__Zerver-REST__ 是一个[golang](http://golang.org)实现的高性能, 简单 RESTFul-api框架 

这是一个[Zerver](http://github.com/cosiner/zerver)的简化版本, 去除了模板, session等功能, 专注与RESTFul风格的api输出, 架构更加简洁明了.

[文档](godoc.org/github.com/cosiner/zerver_rest)

#### 安装
`go get github.com/cosiner/zever_rest`

#### 起步
```Go
package main

import zerver "github.com/cosiner/zever_rest"

func main() {
    server := zever.NewServer()
    server.AddFuncHandler("/", "GET", req zever.Request, resp zever.Response) {
        resp.Write([]byte("Hello World!"))    
    })
    server.Start(":8080")
}
```

-------------------------------------------------------------------------------
#### 特性
* RESTFul风格的路由设计
* 基于前缀树的路由匹配, 放弃正则表达式
* 支持过滤器链
* 内置WebSocket支持

-------------------------------------------------------------------------------
#### 注意
Zerver-REST仅定义了一个编程框架, 核心只有一个路由分发器, 任何其他功能比如session, 模板, 多种输出格式都被移除了, 如果确实需要, 可以尝试[Zerver](http://github.com/cosiner/zerver), 而且并不与golang标准库兼容

##### Router
Zerver-REST's Router is based on prefix-tree, it's transparent to user.
It now support three routes:
* static route: such as /user/12345/info
* variables route: such as /user/:id/info, id is a variable in route, variable value will be injected to Request
* catch-all route: such as /home/*subpath, subpath is a variable in route, it will catch remains all url path, it should appeared in last of route for any route path after it will be ingored

__Route Match__:
The standard Router will not performed backtracking, match priority:static > variable > catchall.

example: route:/user/ab/234, /user/:id/123, url path:/user/ab/123
first, route /user/ab will first matched, but 234 don't match 123, so no route will be matched, regardless /user/:id/123 is actually matched, to fix it,
you should add a route /user/ab/123 as a special case for route /user/:id/123.

__Handler/WebSocketHandler/Filter Match__:
handler, websockethandler is perform full-matched, filter is perform prefix-matched and don't distinguish /use, /user, for /user/123, both /use and /user is matched for filters.

#### URLVarIndexer
If a incoming request is matched by router, variable values exist in matched route will be extracted and packaged as URLVarIndexer.
URLVar(string) return value of url variable
ScanURLVars(...*string) scan values to given addresses, if want to skip a variable, place a 'nil' to it's position
URLVars() []string return all variable values

#### Request/Response
Zerver's Request wrapped standard http.Request, and Response wrapped standard
http.ResponseWriter. The Request is also a URLVarIndexer.

#### [WebSocket]Handler
The most important component of Zerver is Handler, which handles http request for matched route. For HTTP Handler, it accept Request and Response, all things is done with them, Request is input, and Response is output. For WebSocketHandler, it accept WebSocketConn which represent a websocket connection, it's input and output.

#### Filter
Filter will be called before any Handler, and it's only filter Handler, not WebSocketHandler, it accept Request, Response, FilterChain. all matched filters will be packed as a FilterChain, to continue process http request, you must call FilterChain.Filter(Request, Response), otherwise, process is stoped.

Execute sequence of filters is early route first, in one route, it's early added first.






