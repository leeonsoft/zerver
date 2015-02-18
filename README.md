#### Zerver-REST
__Zerver-REST__ is a high-performance, simple, restful api framework for [golang](http://golang.org)

It's a simplified version of [Zever](http://github.com/cosiner/zever), remove template, session, and some other components to get a more clear architecture.

It's mainly designed for restful api service, but you can also use it as a web framework.

-------------------------------------------------------------------------------
#### Install
`go get github.com/cosiner/zever_rest`

-------------------------------------------------------------------------------
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

-------------------------------------------------------------------------------
#### Note
Zerver-REST only define some specifications, it's only a router, and only support output bytes, multiple output format like json/xml/gob should use third party libraries.

-------------------------------------------------------------------------------
##### Router
Zerver-REST's Router is based on prefix-tree.
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

-------------------------------------------------------------------------------
#### Filter
Filter will be called before any Handler, and it's only filter Handler, not WebSocketHandler, in Zerver/Zerver-REST, filters will be packed as a FilterChain, to continue process http request, you must call FilterChain.Filter(Request, Response), otherwise, process is stoped.

Execute sequence of filters is early route first, in one route, it's early added first.

