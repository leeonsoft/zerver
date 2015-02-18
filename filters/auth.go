package filters

import (
	"net/http"

	zerver "github.com/cosiner/zerver_rest"
)

type (
	TokenValidator interface {
		// check whether a token is saved
		IsExist(tok string) bool
		// check whether a token is valid
		IsValid(tok string) bool
	}
)

var (
	TokValidator TokenValidator
	AuthKey      []byte
)

const (
	_AUTH_HEADER = "Authentication"
	_AUTHFAILED  = http.StatusUnauthorized
)

func AuthFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	tok := req.Header(_AUTH_HEADER)
	if tok != "" && TokValidator.IsValid(tok) && TokValidator.IsExist(tok) {
		chain.Filter(req, resp)
	} else {
		resp.ReportStatus(_AUTHFAILED)
	}
}
