package filters

import (
	zerver "github.com/cosiner/zerver_rest"
)

var logger func(v ...interface{})

func NewLogFilter(l func(v ...interface{})) zerver.Filter {
	logger = l
	return (zerver.FilterFunc)(logFilter)
}

func logFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	logger(req.RemoteAddr(), req.URL().String(), req.UserAgent())
	chain(req, resp)
}
