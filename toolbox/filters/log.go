package filters

import (
	"github.com/cosiner/zerver"
)

type LogFilter func(v ...interface{})

func (l LogFilter) Init(*zerver.Server) error { return nil }

func (l LogFilter) Filter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	chain(req, resp)
	l("[AccessLog]", resp.Status(), req.URL().String(), req.RemoteAddr(), req.UserAgent())
}

func (l LogFilter) Destroy() {}
