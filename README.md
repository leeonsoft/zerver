# Zerver
__Zerver__ is a simple, scalable, restful api framework for [golang](http://golang.org).

It's mainly designed for restful api service, without session, template, cookie support, etc.. But you can still use it as a web framework by easily hack it. Documentation can be found at [godoc.org](godoc.org/github.com/cosiner/zerver)

### API Reference
Each file contains a component, all api about this component is defined there.

### Install
`go get github.com/cosiner/zerver`

### Getting Started
```Go
package main

import "github.com/cosiner/zerver"

func main() {
    server := zerver.NewServer()
    server.Get("/", func(req zerver.Request, resp zerver.Response) {
        resp.Write([]byte("Hello World!"))    
    })
    server.Start(":8080")
}
```

### Features
* RESTFul Route
* Tree-based mux/router, support route group
* Filter Chain supported
* Interceptor supported
* Builtin WebSocket supported
* Builtin Task supported

# Exampels
* url variables
```Go
server.Get("/user/info/:id", func(req zerver.Request, resp zerver.Response) {
    resp.Write([]byte("Hello, " + req.URLVar("id")))    
})
server.Get("/home/*subpath", func(req zerver.Request, resp zerver.Response) {
    resp.Write([]byte("You access" + req.URLVar("subpath")))    
})
```

* filter
```Go
type logger func(v ...interface{})
func (l logger) Init(*zerver.Server) error {return nil}
func (l logger) Destroy() {}
func (l logger) Filter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
    l(req.RemoteIP(), req.UserAgent(), req.URL())
    chain(req, resp)
}
server.AddFilter("/", logger(fmt.Println))
```

* interceptor
```Go
func BasicAuthFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
    auth := req.Header("Authorization")
    if auth == "" {
        resp.ReportUnAuthorized()
        return
    }
    req.SetAttr("auth", auth) // do before
    chain(req, resp)
}

func AuthLogFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
    chain(req, resp)
    fmt.Println("Auth success: ", resp.Value()) // do after
}

func AuthHandler(req zerver.Request, resp zerver.Response) {
    resp.Write([]byte(req.Attr("auth")))
    resp.SetValue(true)
}

server.AddFuncHandler("/auth", zerver.InterceptHandler(
    AuthHandler, 
    zerver.FilterFunc(BasicAuthFilter), 
    zerver.FilterFunc(AuthLogFilter),
))
```

More detail please see [wiki page](https://github.com/cosiner/zerver/wiki).

### License
MIT.
