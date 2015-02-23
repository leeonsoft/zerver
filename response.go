package zerver_rest

import (
	"bufio"
	"io"
	"net"
	"net/http"

	. "github.com/cosiner/golib/errors"
)

type (
	Response interface {
		SetHeader(name, value string)
		AddHeader(name, value string)
		RemoveHeader(name string)
		SetContentEncoding(enc string)
		SetContentType(typ string)
		SetCookie(name, value string, lifetime int)
		SetSecureCookie(name, value string, lifetime int)
		DeleteClientCookie(name string)
		ReportStatus(statusCode int)
		Hijack() (net.Conn, *bufio.ReadWriter, error)
		Flush()
		io.Writer
	}

	// response represent a response of request to user
	response struct {
		http.ResponseWriter
		header http.Header
	}

	// marshalFunc is the marshal function type
	marshalFunc func(interface{}) ([]byte, error)
)

// newResponse create a new response, and set default content type to HTML
func newResponse(w http.ResponseWriter) Response {
	resp := &response{
		ResponseWriter: w,
		header:         w.Header(),
	}
	return resp
}

// SetHeader setup response header
func (resp *response) SetHeader(name, value string) {
	resp.header.Set(name, value)
}

// AddHeader add a value to response header
func (resp *response) AddHeader(name, value string) {
	resp.header.Add(name, value)
}

// RemoveHeader remove response header by name
func (resp *response) RemoveHeader(name string) {
	resp.header.Del(name)
}

// contentType return current content type of response
func (resp *response) contentType() string {
	return resp.header.Get(HEADER_CONTENTTYPE)
}

// newCookie create a new Cookie and return it's displayed string
// parameter lifetime is time by second
func (*response) newCookie(name, value string, lifetime int) string {
	return (&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: lifetime,
	}).String()
}

// ReportStatus report an http status with given status code
func (resp *response) ReportStatus(statusCode int) {
	resp.WriteHeader(statusCode)
}

// Hijack hijack response connection
func (resp *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, is := resp.ResponseWriter.(http.Hijacker); is {
		return hijacker.Hijack()
	}
	return nil, nil, Err("Connection not support hijack")
}

// Flush flush response's output
func (resp *response) Flush() {
	if flusher, is := resp.ResponseWriter.(http.Flusher); is {
		flusher.Flush()
	}
}

// SetContentType set content type of response
func (resp *response) SetContentType(typ string) {
	resp.SetHeader(HEADER_CONTENTTYPE, typ)
}

// SetContentEncoding set content encoding of response
func (resp *response) SetContentEncoding(enc string) {
	resp.SetHeader(HEADER_CONTENTENCODING, enc)
}

// SetSecureCookie setup response cookie
func (resp *response) SetCookie(name, value string, lifetime int) {
	resp.AddHeader(HEADER_SETCOOKIE, resp.newCookie(name, value, lifetime))
}

// SetSecureCookie setup response cookie with secureity
func (resp *response) SetSecureCookie(name, value string, lifetime int) {
	resp.SetCookie(name, value, lifetime)
}

// DeleteClientCookie delete user briwser's cookie by name
func (resp *response) DeleteClientCookie(name string) {
	resp.SetCookie(name, "", -1)
}
