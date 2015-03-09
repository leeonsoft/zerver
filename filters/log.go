package filters

import (
	zerver "github.com/cosiner/zerver_rest"
)

type InfoLogger interface {
	Infoln(v ...interface{})
}

var Logger InfoLogger

func LogFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	Logger.Infoln(req.RemoteAddr(), req.URL().String(), req.UserAgent())
	chain(req, resp)
}
