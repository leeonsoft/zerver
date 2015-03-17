package filters

import (
	"net/http"

	"github.com/cosiner/zerver"
	jwt "github.com/cosiner/zerver_jwt"
)

var (
	_jwt          *jwt.JWT
	_authAttrName string
)

const (
	_AUTH_HEADER = "Authorization"
	_AUTHFAILED  = http.StatusUnauthorized
)

func NewJWTAuthFilter(j *jwt.JWT, requestAttrName string) zerver.Filter {
	_jwt = j
	if requestAttrName == "" {
		requestAttrName = "Auth"
	}
	_authAttrName = requestAttrName
	return zerver.FilterFunc(jwtAuthFilter)
}

func jwtAuthFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	if tokstr := req.Header(_AUTH_HEADER); tokstr != "" {
		if tok, err := _jwt.Parse(tokstr); err == nil {
			req.SetAttr(_authAttrName, tok)
			chain(req, resp)
			return
		}
	}
	resp.ReportStatus(_AUTHFAILED)
}
