package zerver

import (
	"github.com/cosiner/golib/crypto"
	. "github.com/cosiner/golib/errors"
)

type (
	// SecureCookieFunc is function type of encode/decode secure cookie
	SecureCookieFunc func(string) string
	// SecureCookie enable secure cookie
	SecureCookie interface {
		// EncodeSecureCookie encode a cookie to secure cookie
		EncodeSecureCookie(string) string
		// DecodeSecureCookie decode a secure cookie to normal cookie
		DecodeSecureCookie(string) string
	}
	// emptySecureCookie is a emptySecureCookie processor
	emptySecureCookie struct{}
	// hmacSecureCookie use hmac to encode/verify a cookie value
	hmacSecureCookie struct {
		secret []byte
	}
)

// NewSecureCookieProcessor create a secure cookie processor, if enable, secret
// must not be empty
func NewSecureCookieProcessor(enable bool, secret []byte) SecureCookie {
	if enable {
		Assert(len(secret) != 0, Err("Secret key for secure cookie can't be empty"))
		return &hmacSecureCookie{secret}
	}
	return emptyCompressor{}
}

func (esc emptyCompressor) EncodeSecureCookie(value string) string { return value }
func (esc emptyCompressor) DecodeSecureCookie(value string) string { return value }

func (hsc *hmacSecureCookie) EncodeSecureCookie(value string) string {
	return crypto.SignSecret(hsc.secret, value)
}

func (hsc *hmacSecureCookie) DecodeSecureCookie(value string) string {
	return crypto.VerifySecret(hsc.secret, value)
}
