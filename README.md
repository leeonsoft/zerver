# Zerver
__Zerver__ is a high-performance, simple, restful api framework for [golang](http://golang.org).

It's mainly designed for restful api service, but you can also use it as a web framework. Documentation can be found at [godoc.org](godoc.org/github.com/cosiner/zerver)

### Install
`go get github.com/cosiner/zever`

### Getting Started
```Go
package main

import "github.com/cosiner/zever"

func main() {
    server := zever.NewServer()
    server.Get("/", func(req zever.Request, resp zever.Response) {
        resp.Write([]byte("Hello World!"))    
    })
    server.Start(":8080")
}
```

# Features
* RESTFul Route
* Tree-based mux/router
* Filter Chain supported
* Interceptor supported
* Builtin WebSocket supported
* Builtin Task supported
* Object Pool

### Note
The core of Zerver only define some specifications, it's only a router, and only support output bytes, multiple output format like json/xml/gob should use third party libraries. And also, the `Request` interface only provide method to get parameters for GET request.

### Router
Zerver's Router is based on prefix-tree, it's transparent to user, standard router only match the path of url.
It now support three routes:
* static route: such as /user/12345/info
* variables route: such as /user/:id/info, id is a variable in route, variable value will be injected to Request
* catch-all route: such as /home/*subpath, subpath is a variable in route, it will catch remains all url path, it should appeared in last of route for any route path after it will be ingored.

The default url path variable count `PathVarCount` is 2, it's used as capacity of create slice to store matched values, it should be manually change to the best dependes on your routes.

__Customization__:
Zerver also support customed Router by `NewServerWith`, user can wrap standard router to get a more powerful Router with great features, such as multiple domains, for convenience, Zerver already provide a `ResponseWrapper` for common purpose and a `HostRouter` for multiple domain.

__Route Match__:
The standard Router will not performed backtracking, match priority:static > variable > catchall.

example: route:/user/ab/234, /user/:id/123, url path:/user/ab/123
first, route /user/ab will first matched, but 234 don't match 123, so no route will be matched, regardless /user/:id/123 is actually matched, to fix it,
you should add a route /user/ab/123 as a special case for route /user/:id/123.

NOTE: because there is only one route tree, and the Router default match handler and filters simultaneously, so the filter route will hide the handler route if it's more accurate by match sequence.

example: a filter match /user/1:id, a handler match /user/:id, if the url path is /user/123, only filter will be matched, the same as /user/:id and /user/*id, only same route pattern can avoid this condition, but normally, it's not necessary to consider this for these routes are rarely appeared.

solution: define a same route pattern as filter for handler

__Handler/WebSocketHandler/Filter Match__:
handler, websockethandler is perform full-matched, filter is perform prefix-matched and don't distinguish /use, /user, for /user/123, both /use and /user is matched for filters.

### URLVarIndexer
If a incoming request is matched by router, variable values exist in matched route will be extracted and packaged as URLVarIndexer.
* `URLVar(string)` return value of url variable
* `ScanURLVars(...*string)` scan values to given addresses, if want to skip a variable, place a 'nil' to it's position
* `URLVars() []string` return all variable values

### Request/Response
Zerver's `Request` wrapped standard `http.Request`, and `Response` wrapped standard `http.ResponseWriter`. The Request is also a `URLVarIndexer`.

### [WebSocket]Handler
The most important component of Zerver is Handler, which handles http request for matched route. For HTTP Handler, it accept Request and Response, all things is done with them, Request is input, and Response is output. For WebSocketHandler, it accept WebSocketConn which represent a websocket connection, it's input and output.

### TaskHandler
Sometimes we need start additional task in `[WebSocket]Handler`, such as send a verify email to user after process user register.

Zerver provide a `TaskHandler` to reach this. `TaskHandler` isn't exposed to client, it's more likely a backend service or service dispatcher.

Register a TaskHandler by `server.AddTaskHandler/AddFuncTaskHandler`, start a task by `server.StartTask/StartAsyncTask`.

### Filter
Filter will be called before any Handler, and it's only filter Handler, not WebSocketHandler, it accept Request, Response, FilterChain. all matched filters will be packed as a FilterChain, to continue process http request, you must call FilterChain.Filter(Request, Response), otherwise, process is stoped.

Execute sequence of filters is early route first, in one route, it's early added first.

### Interceptor
Zerver's interceptor is based on Filter, you can wrap a `HandlerFunc`/`FilterChain` with a series of Filters by `InterceptHandler` to get a more powerful, functions seperated `HandlerFunc`/`FilterChain`

### ObjectPool
Zerver provide a object pool `ServerPool`, which is based on `sync.Pool`.
* `RegisterPool` register a pool.
* `NewFrom` get object from registed pool
* `RecycleTo` recycle object to registed pool

After gc, all things stored in pool will be cleared

### Temporary data store
Zerver provide a temporary data store for data initial before server start, it will be destroyed after server start
