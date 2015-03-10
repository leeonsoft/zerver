package zerver_rest

import "net/url"

type (
	// FilterFunc represent common filter function type,
	FilterFunc func(Request, Response, FilterChain)

	// Filter is an filter that run before or after handler,
	// to modify or check request and response
	// it will be inited on server started, destroyed on server stopped
	Filter interface {
		Init(*Server) error
		Destroy()
		Filter(Request, Response, FilterChain)
	}

	// FilterChain represent a chain of filter, the last is final handler,
	// to continue the chain, must call chain(Request, Response)
	FilterChain HandlerFunc

	// FilterChain interface {
	// 	Continue(Request, Response)
	// }

	// filterChain holds filters and handler
	// if there is no filters, the final handler will be called
	filterChain struct {
		filters []Filter
		handler HandlerFunc
	}

	// handlerInterceptor is a interceptor
	handlerInterceptor struct {
		filter  Filter
		handler HandlerFunc
	}

	// Only mmatch for url.Path = "/", others are not concerned
	RootFilters interface {
		// Filters return all root filters
		Filters(url *url.URL) []Filter
		// AddRootFilter add root filter for "/"
		AddRootFilter(Filter)
	}

	rootFilters []Filter
)

// FilterFunc is a function Filter
func (FilterFunc) Init(*Server) error { return nil }
func (FilterFunc) Destroy()           {}
func (fn FilterFunc) Filter(req Request, resp Response, chain FilterChain) {
	fn(req, resp, chain)
}

func NewRootFilters() RootFilters {
	rfs := make(rootFilters, 0)
	return &rfs
}

// Filters returl all root filters
func (rfs *rootFilters) Filters(*url.URL) []Filter {
	return *rfs
}

// AddRootFilter add root filter
func (rfs *rootFilters) AddRootFilter(filter Filter) {
	*rfs = append(*rfs, filter)
}

// NewFilterChain create a chain of filter
// this method is setup to public for which condition there only one route
// need filter, if add to global router, it will make router match slower
// this method can help for these condition
//
// NOTICE: FilterChain keeps some states, it should be use only once
// if need unlimit FilterChain, use InterceptHandler replace
// but it takes less memory space, InterceptHandler takes more, make your choice
func NewFilterChain(filters []Filter, handler func(Request, Response)) FilterChain {
	if l := len(filters); l == 0 {
		return handler
	} else if l == 1 {
		return newInterceptHandler(handler, filters[0])
	}
	chain := Pool.newFilterChain()
	chain.filters = filters
	chain.handler = handler
	return chain.continueChain
}

// continueChain call next filter, if there is no next filter,then call final handler
// if there is no chains, filterChain will be recycle to pool
// if chain is not continue to last, chain will not be recycled
func (chain *filterChain) continueChain(req Request, resp Response) {
	filters := chain.filters
	if len(filters) == 0 {
		if h := chain.handler; h != nil {
			chain.handler(req, resp)
		}
		chain.destroy()
	} else {
		filter := filters[0]
		chain.filters = filters[1:]
		filter.Filter(req, resp, chain.continueChain)
	}
}

// destroy destroy all reference hold by filterChain, then recycle it to pool
func (chain *filterChain) destroy() {
	chain.filters = nil
	chain.handler = nil
	Pool.recycleFilterChain(chain)
}

// InterceptHandler will create a permanent HandlerFunc/FilterChain, there is no
// states keeps, it can be used without limitation
// InterceptHandler wrap a handler with some filters as a interceptor
func InterceptHandler(handler func(Request, Response), filters ...Filter) func(Request, Response) {
	if len(filters) == 0 {
		return handler
	}
	return newInterceptHandler(InterceptHandler(handler, filters[1:]...), filters[0])
}

func newInterceptHandler(handler func(Request, Response), filter Filter) func(Request, Response) {
	if handler == nil {
		handler = emptyHandlerFunc
	}
	return (&handlerInterceptor{
		filter:  filter,
		handler: handler,
	}).Handle
}

func (interceptor *handlerInterceptor) Handle(req Request, resp Response) {
	interceptor.filter.Filter(req, resp, FilterChain(interceptor.handler))
}
