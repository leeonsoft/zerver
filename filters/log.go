package filters

import (
	zerver "github.com/cosiner/zerver_rest"
)

type InfoLogger interface {
	Infoln(v ...interface{})
}

var logger InfoLogger

func NewLogFilter(l InfoLogger) zerver.Filter {
	logger = l
	return (zerver.FilterFunc)(logFilter)
}

func logFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	logger.Infoln(req.RemoteAddr(), req.URL().String(), req.UserAgent())
	chain(req, resp)
}
