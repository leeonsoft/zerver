[简体中文](README-zh_CN.md)

#### Zerver-REST
__Zerver-REST__ is a high-performance, simple, restful api framework for [golang](http://golang.org)

It's a simplified version of [Zever](http://github.com/cosiner/zever), remove template, session, and some other components to get a more clear architecture.

It's mainly designed for restful api service, but you can also use it as a web framework.

Documentation can be found at [godoc.org/github.com/cosiner/zerver_rest]

#### Install
`go get github.com/cosiner/zever_rest`

#### Getting Started
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
#### Features
* RESTFul Route
* Tree-based mux/router
* Filter Chain supported
* Builtin WebSocket supported
* Builtin Task supported

-------------------------------------------------------------------------------
#### Note
The core of Zerver-REST only define some specifications, it's only a router, and only support output bytes, multiple output format like json/xml/gob should use third party libraries.

##### Router
Zerver-REST's Router is based on prefix-tree, it's transparent to user, standard router only match the path of url.
It now support three routes:
* static route: such as /user/12345/info
* variables route: such as /user/:id/info, id is a variable in route, variable value will be injected to Request
* catch-all route: such as /home/*subpath, subpath is a variable in route, it will catch remains all url path, it should appeared in last of route for any route path after it will be ingored.

__Customization__:
Zerver also support customed Router by `NewServerWith`, user can wrap standard router to get a more powerful Router with great features, such as multiple domains.

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

#### TaskHandler
Sometimes we need start additional task in `[WebSocket]Handler`, such as send a verify email to user after process user register.

Zerver provide a `TaskHandler` to reach this. `TaskHandler` isn't exposed to client, it's more likely a backend service or service dispatcher.

Register a TaskHandler by `server.AddTaskHandler/AddFuncTaskHandler`, start a task by `server.StartTask/StartAsyncTask`.

#### Filter
Filter will be called before any Handler, and it's only filter Handler, not WebSocketHandler, it accept Request, Response, FilterChain. all matched filters will be packed as a FilterChain, to continue process http request, you must call FilterChain.Filter(Request, Response), otherwise, process is stoped.

Execute sequence of filters is early route first, in one route, it's early added first.






