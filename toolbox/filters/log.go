package filters

import (
	"github.com/cosiner/zerver"
)

type LogFilter func(v ...interface{})

func (l LogFilter) Init(*zerver.Server) error { return nil }

func (l LogFilter) Filter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	l(req.RemoteAddr(), req.URL().String(), req.UserAgent())
	chain(req, resp)
}

func (l LogFilter) Destroy() {}
