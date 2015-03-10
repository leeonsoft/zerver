package zerver_rest

import (
	"log"
	"net/http"
	"net/url"

	websocket "github.com/cosiner/zerver_websocket"

	. "github.com/cosiner/golib/errors"
)

type (
	// Server represent a web server
	Server struct {
		Router
		AttrContainer
		RootFilters
		handlerMatcher func(*url.URL) (Handler, URLVarIndexer, []Filter)
	}

	serverGetter interface {
		Server() *Server
	}
)

var disableFilter bool

// NewServer create a new server
func NewServer() *Server {
	return NewServerWith(NewRouter(), nil)
}

// NewServerWith create a new server with given router and root filters
func NewServerWith(rt Router, filters RootFilters) *Server {
	if filters == nil {
		filters = NewRootFilters()
	}
	return &Server{
		Router:        rt,
		AttrContainer: NewLockedAttrContainer(),
		RootFilters:   filters,
	}
}

// Implement serverGetter
func (s *Server) Server() *Server {
	return s
}

// Start start server
func (s *Server) start() {
	if disableFilter {
		s.handlerMatcher = s.matchHandlerNoFilters
	} else {
		s.handlerMatcher = s.MatchHandlerFilters
	}
	log.Println("Init Handlers and Filters")
	s.Router.Init(func(handler Handler) bool {
		OnErrPanic(handler.Init(s))
		return true
	}, func(filter Filter) bool {
		OnErrPanic(filter.Init(s))
		return true
	}, func(websocketHandler WebSocketHandler) bool {
		OnErrPanic(websocketHandler.Init(s))
		return true
	}, func(taskHandler TaskHandler) bool {
		OnErrPanic(taskHandler.Init(s))
		return true
	})
	log.Println("Server Start")
	// destroy temporary data store
	tmpDestroy()
}

// Start start server as http server
func (s *Server) Start(listenAddr string) error {
	s.start()
	return http.ListenAndServe(listenAddr, s)
}

// StartTLS start server as https server
func (s *Server) StartTLS(listenAddr, certFile, keyFile string) error {
	s.start()
	return http.ListenAndServeTLS(listenAddr, certFile, keyFile, s)
}

// ServHttp serve for http reuest
// find handler and resolve path, find filters, process
func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if websocket.IsWebSocketRequest(request) {
		s.serveWebSocket(w, request)
	} else {
		s.serveHTTP(w, request)
	}
}

// StartTask add a synchronous task
func (s *Server) StartTask(async bool, path string, value interface{}) {
	if async {
		go s.serveTask(path, value)
	} else {
		s.serveTask(path, value)
	}
}

// serveWebSocket serve for websocket protocal
func (s *Server) serveWebSocket(w http.ResponseWriter, request *http.Request) {
	handler, indexer := s.MatchWebSocketHandler(request.URL)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if conn, err := websocket.UpgradeWebsocket(w, request, nil); err == nil {
		handler.Handle(newWebSocketConn(s, conn, indexer))
		indexer.destroy()
		Pool.recycleVarIndexer(indexer)
	}
}

func (s *Server) matchHandlerNoFilters(url *url.URL) (Handler, URLVarIndexer, []Filter) {
	handler, indexer := s.MatchHandler(url)
	return handler, indexer, nil
}

// serveHTTP serve for http protocal
func (s *Server) serveHTTP(w http.ResponseWriter, request *http.Request) {
	url := request.URL
	url.Host = request.Host
	handler, indexer, filters := s.handlerMatcher(url)
	requestEnv := Pool.newRequestEnv()
	req, resp := requestEnv.req.init(s, request, indexer), requestEnv.resp.init(w)

	var handlerFunc HandlerFunc
	if handler == nil {
		// resp.setStatus(http.StatusNotFound)
		// handlerFunc = notFoundHandler
		resp.ReportNotFound()
	} else if handlerFunc = handler.Handler(req.Method()); handlerFunc == nil {
		// resp.setStatus(http.StatusMethodNotAllowed)
		// handlerFunc = methodNotAllowedHandler
		resp.ReportMethodNotAllowed()
	}
	chain := NewFilterChain(filters, handlerFunc)
	if url.Path == "/" {
		chain = NewFilterChain(s.Filters(url), chain)
	}
	chain(req, resp)
	req.destroy()
	resp.destroy()
	indexer.destroy()
	Pool.recycleRequestEnv(requestEnv)
	Pool.recycleVarIndexer(indexer)
	Pool.recycleFilters(filters)
}

// serveTask serve for asynchronous task
func (s *Server) serveTask(path string, value interface{}) {
	u, err := url.Parse(path)
	if err != nil {
		PanicServer(err.Error())
	}
	handler, indexer := s.MatchTaskHandler(u)
	if handler == nil {
		PanicServer("No task handler found for " + path)
	}
	handler.Handle(newTask(s, indexer, value))
	indexer.destroy()
	Pool.recycleVarIndexer(indexer)
}

func (s *Server) AddRootFuncFilter(filter FilterFunc) {
	s.AddRootFilter(filter)
}

// Get register a function handler process GET request for given pattern
func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, GET, handlerFunc)
}

// Post register a function handler process POST request for given pattern
func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, POST, handlerFunc)
}

// Put register a function handler process PUT request for given pattern
func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, PUT, handlerFunc)
}

// Delete register a function handler process DELETE request for given pattern
func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, DELETE, handlerFunc)
}

// Patch register a function handler process PATCH request for given pattern
func (s *Server) Patch(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, PATCH, handlerFunc)
}

// PanicServer create a new goroutine, it force panic whole process
func PanicServer(s string) {
	go panic(s)
}
