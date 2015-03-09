package filters

import (
	"github.com/cosiner/golib/crypto"
	zerver "github.com/cosiner/zerver_rest"
)

var SecretKey []byte

type (
	secureRequest struct {
		zerver.Request
	}

	secureResponse struct {
		zerver.Response
	}
)

func (sr secureRequest) SecureCookie(name string) string {
	return crypto.VerifySecret(SecretKey, sr.Cookie(name))
}

func (sr secureResponse) SetSecureCookie(name, value string, lifetime int) {
	sr.SetCookie(name, crypto.SignSecret(SecretKey, value), lifetime)
}

func SecureCookieFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	chain(secureRequest{req}, secureResponse{resp})
}
