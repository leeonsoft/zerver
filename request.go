package zerver_rest

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cosiner/golib/encoding"
)

type (
	Request interface {
		RemoteAddr() string
		RemoteIP() string
		Param(name string) string
		Params(name string) []string
		UserAgent() string
		URL() *url.URL
		Method() string
		ContentType() string
		AcceptEncodings() string
		Header(name string) string
		Cookie(name string) string
		SecureCookie(name string) string
		Server() *Server
		io.Reader
		URLVarIndexer
		setURLVarIndexer(indexer URLVarIndexer)
		setURL(url *url.URL)
	}

	// request represent an income request
	request struct {
		URLVarIndexer
		encoding.PowerReader
		server  *Server
		request *http.Request
		method  string
		header  http.Header
		params  url.Values
	}
)

// newRequest create a new request
func newRequest(requ *http.Request, varIndexer URLVarIndexer) Request {
	req := &request{
		request:       requ,
		header:        requ.Header,
		URLVarIndexer: varIndexer,
	}
	method := requ.Method
	if m := requ.Header.Get("X-HTTP-Method-Override"); method == POST && m != "" {
		method = m
	}
	req.method = strings.ToUpper(method)
	return req
}

func (req *request) setURLVarIndexer(indexer URLVarIndexer) {
	req.URLVarIndexer = indexer
}

func (req *request) setURL(url *url.URL) {
	req.request.URL = url
}

func (req *request) Read(data []byte) (int, error) {
	return req.request.Body.Read(data)
}

func (req *request) Server() *Server {
	return req.server
}

// Method return method of request
func (req *request) Method() string {
	return req.method
}

// Cookie return cookie value with given name
func (req *request) Cookie(name string) string {
	if c, err := req.request.Cookie(name); err == nil {
		return c.Value
	}
	return ""
}

// SecureCookie return secure cookie, currently it's just call Cookie without
// 'Secure', if need this feture, just put an filter before handler
// and override this method
func (req *request) SecureCookie(name string) (value string) {
	return req.Cookie(name)
}

// RemoteAddr return remote address
func (req *request) RemoteAddr() string {
	return req.request.RemoteAddr
}

func (req *request) RemoteIP() string {
	return strings.Split(req.RemoteAddr(), ":")[1]
}

// Param return request parameter with name
func (req *request) Param(name string) (value string) {
	params := req.Params(name)
	if len(params) > 0 {
		value = params[0]
	}
	return
}

// Params return request parameters with name
func (req *request) Params(name string) []string {
	params, request := req.params, req.request
	if params == nil {
		switch req.method {
		case GET:
			params = request.URL.Query()
		default:
			request.ParseForm()
			params = request.PostForm
		}
		req.params = params
	}
	return params[name]
}

// UserAgent return user's agent identify
func (req *request) UserAgent() string {
	return req.Header(HEADER_USERAGENT)
}

// ContentType extract content type form request header
func (req *request) ContentType() string {
	return req.Header(HEADER_CONTENTTYPE)
}

func (req *request) AcceptEncodings() string {
	return req.Header(HEADER_ACCEPTENCODING)
}

// URL return request url
func (req *request) URL() *url.URL {
	return req.request.URL
}

// Header return header value with name
func (req *request) Header(name string) string {
	return req.header.Get(name)
}
