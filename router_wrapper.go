package zerver_rest

import (
	"io"
	"net/url"

	"github.com/cosiner/golib/sys"
)

type (
	RouterStore interface {
		FindRouter(*url.URL) Router
		Iterate(func(ident string, router Router))
	}

	RouterWrapper struct {
		RouterStore
	}

	HostRouter struct {
		hosts   []string
		routers []Router
		filters []RootFilters
	}
)

func NewHostRouter() RouterWrapper {
	return RouterWrapper{
		RouterStore: new(HostRouter),
	}
}

func (hr *HostRouter) AddRouter(host string, rt Router, rootFilters RootFilters) {
	l := len(hr.hosts) + 1
	hosts, routers, filters := make([]string, l), make([]Router, l), make([]RootFilters, l)
	copy(hosts, hr.hosts)
	copy(routers, hr.routers)
	copy(filters, hr.filters)
	hosts[l], routers[l], filters[l] = host, rt, rootFilters
	hr.hosts, hr.routers, hr.filters = hosts, routers, filters
}

// Implement RouterStore

func (hr *HostRouter) FindRouter(url *url.URL) Router {
	host, hosts := url.Host, hr.hosts
	for i := range hosts {
		if hosts[i] == host {
			return hr.routers[i]
		}
	}
	return nil
}

func (hr *HostRouter) Iterate(fn func(string, Router)) {
	hosts, routers := hr.hosts, hr.routers
	for i := range routers {
		fn(hosts[i], routers[i])
	}
}

// Implement RootFilters

// Filters return all root filters
func (hr *HostRouter) Filters(url *url.URL) []Filter {
	host, hosts := url.Host, hr.hosts
	for i := range hosts {
		if hosts[i] == host {
			return hr.filters[i].Filters(url)
		}
	}
	return nil
}

// AddRootFilter add root filter for "/"
func (hr *HostRouter) AddRootFilter(Filter) {
	panicOnAdd()
}

// Init init handlers and filters, websocket handlers
func (rw RouterWrapper) Init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocket func(WebSocketHandler) bool,
	initTask func(TaskHandler) bool) {
	rw.Iterate(func(_ string, rt Router) {
		rt.Init(initHandler, initFilter, initWebSocket, initTask)
	})
}

// Destroy destroy router, also responsible for destroy all handlers and filters
func (rw RouterWrapper) Destroy() {
	rw.Iterate(func(_ string, rt Router) {
		rt.Destroy()
	})
}

func panicOnAdd() error {
	PanicServer("Don't operate with router wrapper directly")
	return nil
}

// AddFuncHandler add a function handler, method are defined as constant string
func (RouterWrapper) AddFuncHandler(string, string, HandlerFunc) error {
	return panicOnAdd()
}

// AddHandler add a handler
func (RouterWrapper) AddHandler(string, Handler) error {
	return panicOnAdd()
}

// AddFuncFilter add function filter
func (RouterWrapper) AddFuncFilter(string, FilterFunc) error {
	return panicOnAdd()
}

// AddFilter add a filter
func (RouterWrapper) AddFilter(string, Filter) error {
	return panicOnAdd()
}

// AddFuncWebSocketHandler add a websocket functionhandler
func (RouterWrapper) AddFuncWebSocketHandler(string, WebSocketHandlerFunc) error {
	return panicOnAdd()
}

// AddWebSocketHandler add a websocket handler
func (RouterWrapper) AddWebSocketHandler(string, WebSocketHandler) error {
	return panicOnAdd()
}

// AddFuncTaskHandler
func (RouterWrapper) AddFuncTaskHandler(string, TaskHandlerFunc) error {
	return panicOnAdd()
}

// AddTaskHandler
func (RouterWrapper) AddTaskHandler(string, TaskHandler) error {
	return panicOnAdd()
}

// MatchHandlerFilters match given url to find all matched filters and final handler
func (rw RouterWrapper) MatchHandlerFilters(url *url.URL) (handler Handler, indexer URLVarIndexer, filters []Filter) {
	if router := rw.FindRouter(url); router != nil {
		handler, indexer, filters = router.MatchHandlerFilters(url)
	}
	return
}

// MatchHandler match given url to find all matched filters and final handler
func (rw RouterWrapper) MatchHandler(url *url.URL) (handler Handler, indexer URLVarIndexer) {
	if router := rw.FindRouter(url); router != nil {
		handler, indexer = router.MatchHandler(url)
	}
	return
}

// MatchWebSocketHandler match given url to find a matched websocket handler
func (rw RouterWrapper) MatchWebSocketHandler(url *url.URL) (handler WebSocketHandler, indexer URLVarIndexer) {
	if router := rw.FindRouter(url); router != nil {
		handler, indexer = router.MatchWebSocketHandler(url)
	}
	return
}

// MatchTaskHandler
func (rw RouterWrapper) MatchTaskHandler(url *url.URL) (handler TaskHandler, indexer URLVarIndexer) {
	if router := rw.FindRouter(url); router != nil {
		handler, indexer = router.MatchTaskHandler(url)
	}
	return
}

type indentWriter struct {
	io.Writer
}

var indentBytes = []byte("    ")
var indentBytesLen = len(indentBytes)

func (w indentWriter) Write(data []byte) (int, error) {
	c, e := w.Writer.Write(indentBytes)
	if e == nil {
		c1, e1 := w.Writer.Write(data)
		c, e = c+c1, e1
	}
	return c, e
}

func (rw RouterWrapper) PrintRouteTree(w io.Writer) {
	rw.Iterate(func(ident string, rt Router) {
		if _, e := sys.WriteStrln(w, ident); e == nil {
			rt.PrintRouteTree(indentWriter{w})
		}
	})
}
