package zerver

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
		RootFilters RootFilters
	}
	// ServerInitializer is a Object which will automaticlly initialed by server if
	// it's added to server, else it should initialed manually
	ServerInitializer interface {
		Init(s *Server) error
	}

	// serverGetter is a server getter
	serverGetter interface {
		Server() *Server
	}
)

// NewServer create a new server
func NewServer() *Server {
	return NewServerWith(nil, nil)
}

// NewServerWith create a new server with given router and root filters
func NewServerWith(rt Router, filters RootFilters) *Server {
	if filters == nil {
		filters = NewRootFilters()
	}
	if rt == nil {
		rt = NewRouter()
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
	OnErrPanic(s.RootFilters.Init(s))
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

// Destroy is mainly designed to stop server, release all resources
// but it seems ther is no approach to stop a running golang but don't exit process
func (s *Server) Destroy() {
	s.RootFilters.Destroy()
	s.Router.Destroy()
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
		indexer.destroySelf()
	}
}

// serveHTTP serve for http protocal
func (s *Server) serveHTTP(w http.ResponseWriter, request *http.Request) {
	url := request.URL
	url.Host = request.Host
	handler, indexer, filters := s.MatchHandlerFilters(url)
	requestEnv := Pool.newRequestEnv()
	req, resp := requestEnv.req.init(s, request, indexer), requestEnv.resp.init(w)

	var chain FilterChain
	if handler == nil {
		resp.ReportNotFound()
	} else if chain = FilterChain(IndicateHandler(req.Method(), handler)); chain == nil {
		resp.ReportMethodNotAllowed()
	}
	chain = newFilterChain(filters, chain)
	if url.Path == "/" {
		chain = newFilterChain(s.RootFilters.Filters(url), chain)
	}
	chain(req, resp)
	req.destroy()
	resp.destroy()
	Pool.recycleRequestEnv(requestEnv)
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
	task := newTask(s, indexer, value)
	handler.Handle(task)
	task.destroy()
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
