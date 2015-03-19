package filters

import (
	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/zerver"
	jwt "github.com/cosiner/zerver_jwt"
)

const (
	_HEADER_AUTHRIZATION = "Authorization"
)

type JWTAuthFilter struct {
	JWT               *jwt.JWT
	AuthTokenAttrName string
}

func (j *JWTAuthFilter) Init(s *zerver.Server) error {
	if j.JWT == nil {
		return Err("jwt token generator/validator can't be nil")
	}
	if j.JWT.Keyfunc == nil {
		return Err("jwt secret key getter can't be nil")
	}
	if j.AuthTokenAttrName == "" {
		j.AuthTokenAttrName = "AuthToken"
	}
	if j.JWT.SigningMethod == nil {
		j.JWT.SigningMethod = jwt.SigningMethodHS256
	}
	return nil
}

func (j *JWTAuthFilter) Filter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	if tokstr := req.Header(_HEADER_AUTHRIZATION); tokstr != "" {
		if tok, err := j.JWT.Parse(tokstr); err == nil {
			req.SetAttr(j.AuthTokenAttrName, tok)
			chain(req, resp)
			return
		}
	}
	resp.ReportUnauthorized()
}

func (j *JWTAuthFilter) Destroy() {}
