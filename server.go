package zerver

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/cosiner/zerver/websocket"

	. "github.com/cosiner/golib/errors"
)

//==============================================================================
//                           Server Init
//==============================================================================
type (
	// ServerConfig is all config of server
	ServerConfig struct {
		Router             Router             // router
		ErrorHandlers      ErrorHandlers      // error handlers
		SessionDisable     bool               // session disble
		SessionManager     SessionManager     // session manager
		SessionStore       SessionStore       // session store
		StoreConfig        string             // session store config
		SessionLifetime    int64              // session lifetime
		TemplateEngine     TemplateEngine     // template engine
		FilterForward      bool               // fliter forward request
		SecureCookieEnable bool               // enable secure cookie
		SecretKey          []byte             // secret key
		GzipResponse       bool               // compress response use gzip
		XsrfEnable         bool               // enable xsrf cookie
		XsrfFlushInterval  int                // xsrf cookie value flush interval
		XsrfLifetime       int                // xsrf cookie expire
		XsrfTokenGenerator XsrfTokenGenerator // xsrf token generator
	}

	// Server represent a web server
	Server struct {
		AttrContainer
		Router
		SessionManager
		ErrorHandlers
		TemplateEngine
		xsrfProcessor         Xsrf
		responseCompressor    ResponseCompressor
		secureCookieProcessor SecureCookie
		filterForward         bool
		initOnce              *sync.Once
	}
)

// NewServer create a new server
func NewServer() *Server {
	return &Server{
		AttrContainer:  NewLockedAttrContainer(),
		Router:         NewRouter(),
		ErrorHandlers:  NewErrorHandlers(),
		TemplateEngine: NewTemplateEngine(),
		initOnce:       new(sync.Once),
	}
}

// init init server config
func (sc *ServerConfig) init() {
	if !sc.SessionDisable {
		if sc.SessionStore == nil {
			sc.StoreConfig = DEF_SESSION_MEMSTORE_CONF
			sc.SessionStore = NewMemStore()
		}
		if sc.SessionManager == nil {
			sc.SessionManager = NewSessionManager()
		}

		if sc.SessionLifetime == 0 {
			sc.SessionLifetime = DEF_SESSION_LIFETIME
		}
	} else {
		sc.SessionManager = newPanicSessionManager()
	}
}

// Init init server with given config, it do only once, if the config contains
// router, it must be called before any operation of add Handler/Filter/WebsocketHandler
// otherwise, ann added Handler/Filter/WebsocketHandler will be added to the old router
// templates and locales are inited last, so Handler/Filter/WebSockethandler's Init
// can use server to add templates and other resources
func (s *Server) Init(conf *ServerConfig) {
	s.initOnce.Do(func() {
		if conf.Router != nil {
			s.Router = conf.Router
		}
		if conf.ErrorHandlers != nil {
			s.ErrorHandlers = conf.ErrorHandlers
		}
		if conf.TemplateEngine != nil {
			s.TemplateEngine = conf.TemplateEngine
		}
		strach.setServerConfig(conf)
	})
}

// Start start server
func (s *Server) start() {
	srvConf := strach.serverConfig()
	if srvConf == nil {
		srvConf = new(ServerConfig)
	}
	srvConf.init()
	if srvConf.XsrfEnable {
		xsrfFlush := srvConf.XsrfFlushInterval
		if xsrfFlush <= 0 {
			xsrfFlush = XSRF_DEFFLUSH
		}
		xsrfLifetime := srvConf.XsrfLifetime
		if xsrfLifetime <= 0 {
			xsrfLifetime = XSRF_DEFLIFETIME
		}
		log.Println("XsrfConfig: FlushInterval:", xsrfFlush, " CookieLifetime:", xsrfLifetime)
		xsrfGen := srvConf.XsrfTokenGenerator
		s.xsrfProcessor = NewXsrf(xsrfGen, xsrfLifetime)
	} else {
		s.xsrfProcessor = emptyXsrf{}
	}
	log.Println("Enable Response Compress:", srvConf.GzipResponse)
	s.responseCompressor = NewResponseCompressor(srvConf.GzipResponse)
	log.Println("Enable Secure Cookie:", srvConf.SecureCookieEnable)
	s.secureCookieProcessor = NewSecureCookieProcessor(srvConf.SecureCookieEnable, srvConf.SecretKey)
	log.Println("Enable Filters for Forward Request:", srvConf.FilterForward)
	s.filterForward = srvConf.FilterForward

	manager := srvConf.SessionManager
	if !srvConf.SessionDisable {
		log.Println("Init Session Store and Manager")
		store := srvConf.SessionStore
		OnErrPanic(store.Init(srvConf.StoreConfig))
		OnErrPanic(manager.Init(store, srvConf.SessionLifetime))
	}
	s.SessionManager = manager

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
	})

	log.Println("Init Error Handlers")
	s.ErrorHandlers.Init(s)

	log.Println("Compile Templates")
	log.Println("	templates:", strach.tmpls())
	log.Println("	delimeters:", strach.tmplDelimsSlice())
	OnErrPanic(s.CompileTemplates())

	log.Println("Initial I18N Locales")
	localeFiles := strach.localeFiles()
	if len(localeFiles) > 0 {
		for _, localeFile := range localeFiles {
			OnErrPanic(_tr.load(localeFile))
		}
		defaultLocale := strach.defaultLocale()
		if defaultLocale == "" {
			defaultLocale = "en_US"
		}
		OnErrPanic(_tr.setDefaultLocale(defaultLocale))
	}

	strach.destroy()
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

//==============================================================================
//                          Server Process
//==============================================================================
// ServHttp serve for http reuest
// find handler and resolve path, find filters, process
func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if websocket.IsWebSocketRequest(request) {
		s.serveWebSocket(w, request)
	} else {
		s.serveHttp(w, request)
	}
}

// serveWebSocket serve for websocket protocal
func (s *Server) serveWebSocket(w http.ResponseWriter, request *http.Request) {
	handler, indexer := s.MatchWebSocketHandler(request.URL)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if conn, err := websocket.UpgradeWebsocket(w, request, nil); err == nil {
		handler.Handle(newWebSocketConn(conn).setUrlVarIndexer(indexer))
	}
}

// serveHttp serve for http protocal
func (s *Server) serveHttp(w http.ResponseWriter, request *http.Request) {
	ctx := newContext(s, w, request)
	req, resp := newRequest(ctx, request), newResponse(ctx, w)
	ctx.init(req, resp)
	s.processHttpRequest(request.URL, req, resp, false) // process
	ctx.destroy()
}

// processHttpRequest do process http request
func (s *Server) processHttpRequest(url *url.URL, req *request, resp *response, forward bool) {
	var (
		xsrfError   bool
		handlerFunc HandlerFunc
		method      = req.Method()
		matchFunc   = s.MatchHandlerFilters
	)
	if !forward { // check xsrf error
		if method == GET {
			resp.setXsrfToken(s.xsrfProcessor.Set(req, resp))
		} else {
			xsrfError = !s.xsrfProcessor.IsValid(req)
		}
	} else if !s.filterForward {
		matchFunc = s.MatchHandler
	}
	handler, indexer, filters := matchFunc(url)
	handlerFunc = s.indicateHandlerFunc(method, xsrfError, handler)
	wc := s.responseCompressor.EnableCompress(req, resp)
	if len(filters) == 0 {
		handlerFunc(req.setVarIndexer(indexer), resp)
	} else {
		NewFilterChain(filters, handlerFunc).Filter(req.setVarIndexer(indexer), resp)
	}
	s.responseCompressor.CloseResponseCompress(resp, wc)
}

// indicateHandlerFunc indicate handler function from method, has xsrf error and handler
func (s *Server) indicateHandlerFunc(method string, xsrfError bool, handler Handler) (
	handlerFunc HandlerFunc) {
	if handler != nil {
		if handlerFunc = IndicateHandler(method, handler); handlerFunc == nil {
			handlerFunc = s.MethodNotAllowedHandler()
		} else if xsrfError {
			if handler, is := handler.(XsrfErrorHandler); is {
				handlerFunc = handler.HandleXsrfError
			} else {
				handlerFunc = s.XsrfErrorHandler()
			}
		}
	} else { // no handler means no resource there
		handlerFunc = s.NotFoundHandler()
	}
	return
}

//==============================================================================
//                           Server I18N
//==============================================================================
// AddLocale add an locale file or dir
func (*Server) AddLocale(path string) {
	strach.addLocaleFile(path)
}

// SetDefaultLocale set default locale for i18n
func (*Server) SetDefaultLocale(locale string) {
	strach.setDefaultLocale(locale)
}

//==============================================================================
//                           Server Handler
//==============================================================================
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

//==============================================================================
//                           Server util
//==============================================================================
// PanicServer panic server by create a new goroutine then panic
func PanicServer(str string) {
	go panic(str)
}
