package zerver_rest

import (
	"net/http"
)

type (
	// HandlerFunc is the common request handler function type
	HandlerFunc func(Request, Response)

	// Handler is an common interface of request handler
	// it will be inited on server started, destroyed on server stopped
	// and use it's method to process income request
	// if want to custom the relation of method and handler, just embed
	// EmptyHandler to your Handler, then implement interface MethodIndicator
	Handler interface {
		Init(*Server) error
		Destroy()
		Get(Request, Response)    // query
		Post(Request, Response)   // add
		Delete(Request, Response) // delete
		Put(Request, Response)    // create or replace
		Patch(Request, Response)  // update
	}

	// MethodIndicator is an interface for user handler to
	// custom method handle functions
	MethodIndicator interface {
		// Handler return an method handle function by method name
		// if nill returned, means access forbidden
		Handler(method string) HandlerFunc
	}

	// EmptyHandler is an empty handler for user to embed
	EmptyHandler struct{}

	// funcHandler is a handler that use user customed handler function
	// it a funcHandler is defined after a normal handler with same pattern,
	// regardless of only one method is used in funcHandler like Get, it's impossible
	// to use the Post of normal handler, whole normal handler is hidded by this
	// funcHandler, and if other method handler like Post, Put is not set,
	// user access of these method is forbiddened
	funcHandler struct {
		EmptyHandler
		handlers map[string]HandlerFunc
	}
)

// IndicateHandler indicate handler function from a handler and method
func IndicateHandler(method string, handler Handler) HandlerFunc {
	switch handler := handler.(type) {
	case MethodIndicator:
		return handler.Handler(method)
	default:
		return standardIndicate(method, handler)
	}
}

// standardIndicate normal indicate method handle function
// each method indicate the function with same name, such as GET->Get...
func standardIndicate(method string, handler Handler) (handlerFunc HandlerFunc) {
	switch method {
	case GET:
		handlerFunc = handler.Get
	case POST:
		handlerFunc = handler.Post
	case DELETE:
		handlerFunc = handler.Delete
	case PUT:
		handlerFunc = handler.Put
	case PATCH:
		handlerFunc = handler.Patch
	}
	return
}

// newFuncHandler create a new function handler
func newFuncHandler() *funcHandler {
	return &funcHandler{
		handlers: make(map[string]HandlerFunc),
	}
}

// funcHandler implements MethodIndicator interface for custom method handler
func (fh *funcHandler) Handler(method string) (handlerFunc HandlerFunc) {
	return fh.handlers[method]
}

// setMethodHandler setup method handler for funcHandler
func (fh *funcHandler) setMethodHandler(method string, handlerFunc HandlerFunc) {
	fh.handlers[method] = handlerFunc
}

// EmptyHandler methods
func (EmptyHandler) Init(*Server) error              { return nil }
func (EmptyHandler) Destroy()                        {}
func (EmptyHandler) Get(_ Request, resp Response)    { resp.ReportStatus(http.StatusMethodNotAllowed) }
func (EmptyHandler) Post(_ Request, resp Response)   { resp.ReportStatus(http.StatusMethodNotAllowed) }
func (EmptyHandler) Delete(_ Request, resp Response) { resp.ReportStatus(http.StatusMethodNotAllowed) }
func (EmptyHandler) Put(_ Request, resp Response)    { resp.ReportStatus(http.StatusMethodNotAllowed) }
func (EmptyHandler) Patch(_ Request, resp Response)  { resp.ReportStatus(http.StatusMethodNotAllowed) }
