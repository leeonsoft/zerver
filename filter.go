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

	simpleFilterChain struct {
		runFilter bool
		filter    Filter
		handler   HandlerFunc
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

// There three methods to use Filter/FilterFunc
// 1. add filters to router for common filter
//    such as router.AddFilter/AddFuncFilter(pattern, filter)
// 2. add filters to handler for multiple methods handler
//     such as add
// 3. add filters to function handler for intercept

// NewFilterChain create a chain of filter
// this method is setup to public for which condition there only one route
// need filter, if add to global router, it will make router match slower
// this method can help for these condition
func NewFilterChain(filters []Filter, handler func(Request, Response)) FilterChain {
	if handler == nil {
		handler = emptyHandlerFunc
	}
	if l := len(filters); l == 0 {
		return handler
	} else if l == 1 {
		return newSimpleFilterChain(filters[0], handler)
	}
	chain := Pool.newFilterChain()
	chain.filters = filters
	chain.handler = handler
	return chain.continueChain
}

// InterceptHandle wrap a handler with a filter function as a interceptor
func InterceptHandler(filter Filter, handler func(Request, Response)) func(Request, Response) {
	if handler == nil {
		handler = emptyHandlerFunc
	}
	if filter == nil {
		return handler
	}
	return newSimpleFilterChain(filter, handler)
}

func newSimpleFilterChain(filter Filter, handler func(Request, Response)) func(Request, Response) {
	return (&simpleFilterChain{
		filter:  filter,
		handler: handler,
	}).continueChain
}

// continueChain call next filter, if there is no next filter,then call final handler
// if there is no chains, filterChain will be recycle to pool
// if chain is not continue to last, chain will not be recycled
func (chain *filterChain) continueChain(req Request, resp Response) {
	filters := chain.filters
	if len(filters) == 0 {
		chain.destroy()
		chain.handler(req, resp)
	} else {
		filter := filters[0]
		chain.filters = filters[1:]
		filter.Filter(req, resp, chain.continueChain)
	}
}

func (chain *simpleFilterChain) continueChain(req Request, resp Response) {
	runFilter := !chain.runFilter
	chain.runFilter = runFilter
	if runFilter {
		chain.filter.Filter(req, resp, chain.continueChain)
	} else {
		chain.handler(req, resp)
	}
}

// destroy destroy all reference hold by filterChain, then recycle it to pool
func (chain *filterChain) destroy() {
	chain.filters = nil
	chain.handler = nil
	Pool.recycleFilterChain(chain)
}
