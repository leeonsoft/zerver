package zerver_rest

type (
	// FilterFunc represent common filter function type,
	FilterFunc func(Request, Response, FilterChain)
	// FilterChain represent a chain of filter, the last is final handler,
	// to continue the chain, must call chain(Request, Response)
	FilterChain func(Request, Response)

	// Filter is an filter that run before or after handler,
	// to modify or check request and response
	// it will be inited on server started, destroyed on server stopped
	Filter interface {
		Init(*Server) error
		Destroy()
		Filter(Request, Response, FilterChain)
	}

	// FilterChain interface {
	// 	Continue(Request, Response)
	// }

	// filterChain holds filters and handler
	// if there is no filters, the final handler will be called
	filterChain struct {
		index   int
		filters []Filter
		handler HandlerFunc
	}
)

// FilterFunc is a function Filter
func (FilterFunc) Init(*Server) error { return nil }
func (FilterFunc) Destroy()           {}
func (fn FilterFunc) Filter(req Request, resp Response, chain FilterChain) {
	fn(req, resp, chain)
}

// NewFilterChain create a chain of filter
// this method is setup to public for which condition there only one route
// need filter, if add to global router, it will make router match slower
// this method can help for these condition
func NewFilterChain(filters []Filter, handler func(Request, Response)) FilterChain {
	if len(filters) == 0 {
		return handler
	}
	chain := Pool.newFilterChain()
	chain.index = -1
	chain.filters = filters
	chain.handler = handler
	return chain.continueChain
}

// Filter call next filter, if there is no next filter,then call final handler
// if there is no chains, filterChain will be recycle to pool
// if chain is not continue to last, chain will not be recycled
func (chain *filterChain) continueChain(req Request, resp Response) {
	chain.index++
	index, filters := chain.index, chain.filters
	if index == len(filters) {
		if handler := chain.handler; handler != nil {
			chain.destroy()
			handler(req, resp)
		}
	} else {
		filters[index].Filter(req, resp, chain.continueChain)
	}
}

// destroy destroy all reference hold by filterChain, then recycle it to pool
func (chain *filterChain) destroy() {
	chain.filters = nil
	chain.handler = nil
	Pool.recycleFilterChain(chain)
}
