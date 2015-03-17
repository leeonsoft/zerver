package filters

import (
	"net/http"

	"github.com/cosiner/zerver"
	jwt "github.com/cosiner/zerver_jwt"
)

var (
	JWT               *jwt.JWT
	AuthTokenAttrName string
)

const (
	_HEADER_AUTHRIZATION = "Authorization"
	_AUTHFAILED          = http.StatusUnauthorized
)

func JWTAuthFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	if tokstr := req.Header(_HEADER_AUTHRIZATION); tokstr != "" {
		if tok, err := JWT.Parse(tokstr); err == nil {
			req.SetAttr(AuthTokenAttrName, tok)
			chain(req, resp)
			return
		}
	}
	resp.ReportStatus(_AUTHFAILED)
}
