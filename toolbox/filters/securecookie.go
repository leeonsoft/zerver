package filters

import (
	"github.com/cosiner/golib/crypto"
	zerver "github.com/cosiner/zerver_rest"
)

var secretKey []byte

type (
	secureRequest struct {
		zerver.Request
	}

	secureResponse struct {
		zerver.Response
	}
)

func (sr secureRequest) SecureCookie(name string) string {
	return crypto.VerifySecret(secretKey, sr.Cookie(name))
}

func (sr secureResponse) SetSecureCookie(name, value string, lifetime int) {
	sr.SetCookie(name, crypto.SignSecret(secretKey, value), lifetime)
}

func NewSecureCookieFilter(secret []byte) zerver.Filter {
	secretKey = secret
	return (zerver.FilterFunc)(secureCookieFilter)
}

func secureCookieFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	chain(secureRequest{req}, secureResponse{resp})
}
