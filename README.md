#### Zerver
Zerver is a high-performance, component-oriented, restful web framework for [golang](http://golang.org)

#### Install
`go get github.com/cosiner/zever`

#### Getting Started
```Go
package main

import "github.com/cosiner/zever"

func main() {
    server := zever.NewServer()
    server.AddFuncHandler("/", "GET", req zever.Request, resp zever.Response) {
        resp.WriteString("Hello World!")    
    })
    server.Start("", "8080")
}
```

#### Features
* RESTful route is well supported
* Tree-based mux/router: route is arranged as a prefix-tree, and the performance is similar to [httprouter](https://github.com/julienschmidt/httprouter), more benchmark is needed
* Component-Oriented Design of Server: almost each function of server is replaceable component

#### Note
Zever is not compatiable with standard http.HandlerFunc interface, instead of,
zever design it's own Handler, Request, Response, and even Filter


#### Server
In Zever, Server is only a request dispatcher which implements standard http.Handler interface, any other function is complemeted with other component managed by it.
There seven components in Server:
* AttrContainer: store sever-scoped attributes
* Router: responsible for add route and route finding, also extract url variable values
* SessionManager: manage server session
* ErrorHandlers: process http error, such as page not found, access forbidden, etc..
* TemplateEngine: add/render templates
* ResponseCompressor: compress response output
* SecureCookie: process with secure cookie

Components is all designed as interface, you can replace each component with your owns, but you should comply with some rules to make these components work well.

#### ServerConfig
ServerConfig represent all configureable attributes of Server, it only used only once for Server initial when call Server.Init.
* Router: in Zever, all routes is managed by Router, so Router should be configured before any operation of Add Handler/Filter/WebSocketHandler
* ErrorHandlers
* SessionDisable: if disable session, any calls of request.Session() will cause panic of server
* SessionManager: SessionManager, if empty, Zever will use default session manager
* SessionStore: session store, if empty, Zever will use default memory store
* StoreConfig: session store config
* SessionLifetime: lifetime of client session cookie
* TemplateEngine: template engine, if empty, Zever will use standard html.Template
* FilterForward: whether match filters when request is forward from another url
* SecureCookieEnable: whether enable secure cookie
* SecretKey: secret key to create secure cookie
* GzipResponse: whether enable gzip response compress
* XsrfEnable: whether enable xsrf checking
* XsrfFlushInterval: xsrf token value flush interval
* XsrfLifetime: lifetime of client xsrf token cookie
* XsrfTokenGenerator: xsrf token generator, if empty, use default random generator

