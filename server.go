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
	}

	serverGetter interface {
		Server() *Server
	}
)

// NewServer create a new server
func NewServer() *Server {
	return &Server{Router: NewRouter(), AttrContainer: NewLockedAttrContainer()}
}

// NewServerWith create a new server with given router
func NewServerWith(rt Router) *Server {
	return &Server{Router: rt, AttrContainer: NewLockedAttrContainer()}
}

func (s *Server) Server() *Server {
	return s
}

// Start start server
func (s *Server) start() {
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
func (s *Server) StartTask(path string, value interface{}) {
	s.serveTask(path, value)
}

// StartAsyncTask add a asynchronous task
func (s *Server) StartAsyncTask(path string, value interface{}) {
	go s.serveTask(path, value)
}

// serveWebSocket serve for websocket protocal
func (s *Server) serveWebSocket(w http.ResponseWriter, request *http.Request) {
	handler, indexer := s.MatchWebSocketHandler(request.URL)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if conn, err := websocket.UpgradeWebsocket(w, request, nil); err == nil {
		handler.Handle(newWebSocketConn(s, conn, indexer))
	}
}

// serveHTTP serve for http protocal
func (s *Server) serveHTTP(w http.ResponseWriter, request *http.Request) {
	url := request.URL
	url.Host = request.Host
	handler, indexer, filters := s.MatchHandlerFilters(url)
	req, resp := newRequest(s, request, indexer), newResponse(w)
	if handler != nil {
		if handlerFunc := IndicateHandler(req.Method(), handler); handlerFunc == nil {
			resp.ReportStatus(http.StatusMethodNotAllowed)
		} else if len(filters) == 0 {
			handlerFunc(req, resp)
		} else {
			NewFilterChain(filters, handlerFunc).Filter(req, resp)
		}
	} else { // no handler means no resource there
		resp.ReportStatus(http.StatusNotFound)
	}
	req.destroy()
	resp.destroy()
}

// serveTask serve for asynchronous task
func (s *Server) serveTask(path string, value interface{}) {
	u, err := url.Parse(path)
	if err != nil {
		s.PanicServer(err.Error())
	}
	handler, indexer := s.MatchTaskHandler(u)
	if handler == nil {
		s.PanicServer("No task handler found for " + path)
	}
	handler.Handle(newTask(s, indexer, value))
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
func (*Server) PanicServer(s string) {
	go panic(s)
}
