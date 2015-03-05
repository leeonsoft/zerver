package zerver_rest

import "net/url"

type (
	RouterStore interface {
		AddRouter(ident string, rt Router)
		FindRouter(*url.URL) Router
		Iterate(func(Router))
	}

	RouterWrapper struct {
		RouterStore
	}

	HostRouter struct {
		hosts   []string
		routers []Router
	}
)

func NewHostRouter() RouterWrapper {
	return RouterWrapper{
		RouterStore: new(HostRouter),
	}
}

// Init init handlers and filters, websocket handlers
func (rw RouterWrapper) Init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocket func(WebSocketHandler) bool,
	initTask func(TaskHandler) bool) {
	rw.Iterate(func(rt Router) {
		rt.Init(initHandler, initFilter, initWebSocket, initTask)
	})
}

// Destroy destroy router, also responsible for destroy all handlers and filters
func (rw RouterWrapper) Destroy() {
	rw.Iterate(func(rt Router) {
		rt.Destroy()
	})
}

// AddFuncHandler add a function handler, method are defined as constant string
func (rw RouterWrapper) AddFuncHandler(pattern string, method string, handler HandlerFunc) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddHandler add a handler
func (rw RouterWrapper) AddHandler(pattern string, handler Handler) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddFuncFilter add function filter
func (rw RouterWrapper) AddFuncFilter(pattern string, filter FilterFunc) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddFilter add a filter
func (rw RouterWrapper) AddFilter(pattern string, filter Filter) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddFuncWebSocketHandler add a websocket functionhandler
func (rw RouterWrapper) AddFuncWebSocketHandler(pattern string, handler WebSocketHandlerFunc) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddWebSocketHandler add a websocket handler
func (rw RouterWrapper) AddWebSocketHandler(pattern string, handler WebSocketHandler) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddFuncTaskHandler
func (rw RouterWrapper) AddFuncTaskHandler(pattern string, handler TaskHandlerFunc) error {
	panic("Don't operate with router wrapper directly")
	return nil
}

// AddTaskHandler
func (rw RouterWrapper) AddTaskHandler(pattern string, handler TaskHandler) error {
	panic("Don't operate with router wrapper directly")
	return nil
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

func (hr *HostRouter) AddRouter(host string, rt Router) {
	l := len(hr.hosts) + 1
	hosts, routers := make([]string, l), make([]Router, l)
	copy(hosts, hr.hosts)
	copy(routers, hr.routers)
	hosts[l], routers[l] = host, rt
	hr.hosts, hr.routers = hosts, routers
}

func (hr *HostRouter) FindRouter(url *url.URL) Router {
	host, hosts := url.Host, hr.hosts
	for i := range hosts {
		if hosts[i] == host {
			return hr.routers[i]
		}
	}
	return nil
}

func (hr *HostRouter) Iterate(fn func(Router)) {
	routers := hr.routers
	for i := range routers {
		fn(routers[i])
	}
}
