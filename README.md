#### Zerver
__Zerver__ is a high-performance, component-oriented, restful web framework for [golang](http://golang.org)

-------------------------------------------------------------------------------
#### Install
`go get github.com/cosiner/zever`

-------------------------------------------------------------------------------
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

-------------------------------------------------------------------------------
#### Features
* __RESTful__ route is well supported
* __Tree-based mux/router__: route is arranged as a prefix-tree, and the performance is similar to [httprouter](https://github.com/julienschmidt/httprouter), more benchmark is needed
* __Component-Oriented Design__ of Server: almost each function of server is replaceable component

-------------------------------------------------------------------------------
#### Note
Zever is not compatiable with standard http.HandlerFunc interface, instead of,
zever design it's own Handler, Request, Response, and even Filter

-------------------------------------------------------------------------------
#### Server
In Zever, Server is only a request dispatcher which implements standard http.Handler interface, any other function is completed with other component managed by it.

There seven components in Server:
* __AttrContainer__: store sever-scoped attributes
* __Router__: responsible for add route and route finding, also extract url variable values
* __SessionManager__: manage server session
* __ErrorHandlers__: process http error, such as page not found, access forbidden, etc..
* __TemplateEngine__: add/render templates
* __ResponseCompressor__: compress response output
* __SecureCookie__: process with secure cookie

Components is all designed as interface, you can replace each component with your owns, but you should comply with some rules to make these components work well.

-------------------------------------------------------------------------------
#### ServerConfig
ServerConfig represent all configureable attributes of Server, it only used only once for Server initial when call Server.Init.
* __Router__: in Zever, all routes is managed by Router, so Router should be configured before any operation of Add Handler/Filter/WebSocketHandler
* __ErrorHandlers__
* __SessionDisable__: if disable session, any calls of request.Session() will cause panic of server
* __SessionManager__: SessionManager, if empty, Zever will use default session manager
* __SessionStore__: session store, if empty, Zever will use default memory store
* __StoreConfig__: session store config
* __SessionLifetime__: lifetime of client session cookie
* __TemplateEngine__: template engine, if empty, Zever will use standard html.Template
* __FilterForward__: whether match filters when request is forward from another url
* __SecureCookieEnable__: whether enable secure cookie
* __SecretKey__: secret key to create secure cookie
* __GzipResponse__: whether enable gzip response compress
* __XsrfEnable__: whether enable xsrf checking
* __XsrfFlushInterval__: xsrf token value flush interval
* __XsrfLifetime__: lifetime of client xsrf token cookie
* __XsrfTokenGenerator__: xsrf token generator, if empty, use default random generator
