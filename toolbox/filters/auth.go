package filters

import (
	"net/http"

	"github.com/cosiner/zerver"
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
	tokValidator TokenValidator
)

const (
	_AUTH_HEADER = "Authentication"
	_AUTHFAILED  = http.StatusUnauthorized
)

func NewAuthFilter(validator TokenValidator) zerver.Filter {
	tokValidator = validator
	return (zerver.FilterFunc)(authFilter)
}

func authFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	tok := req.Header(_AUTH_HEADER)
	if tok != "" && tokValidator.IsValid(tok) && tokValidator.IsExist(tok) {
		chain(req, resp)
	} else {
		resp.ReportStatus(_AUTHFAILED)
	}
}
