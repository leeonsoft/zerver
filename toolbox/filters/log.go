package filters

import (
	"github.com/cosiner/zerver"
)

var Logger func(v ...interface{})

func AccessLogFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	Logger(req.RemoteAddr(), req.URL().String(), req.UserAgent())
	chain(req, resp)
}
