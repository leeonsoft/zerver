package zerver

import (
	"bufio"
	"net"
	"net/http"

	"github.com/cosiner/golib/types"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/golib/encoding"
)

type (
	Response interface {
		Render(tmpl string) error
		RenderWith(tmpl string, value interface{}) error
		SetHeader(name, value string)
		AddHeader(name, value string)
		RemoveHeader(name string)
		SetContentEncoding(enc string)
		SetContentType(typ string)
		SetCookie(name, value string, lifetime int)
		SetSecureCookie(name, value string, lifetime int)
		DeleteClientCookie(name string)
		Redirect(url string)
		PermanentRedirect(url string)
		ReportStatus(statusCode int)
		Hijack() (net.Conn, *bufio.ReadWriter, error)
		Flush()
		AttrContainer
		types.WriterChain
		encoding.CleanPowerWriter
	}

	// response represent a response of request to user
	response struct {
		*context
		types.WriterChain
		encoding.CleanPowerWriter
		header http.Header
	}

	// marshalFunc is the marshal function type
	marshalFunc func(interface{}) ([]byte, error)
)

// newResponse create a new response, and set default content type to HTML
func newResponse(ctx *context, w http.ResponseWriter) *response {
	chain := types.NewWriterChain(w)
	pw := encoding.NewPowerWriter(chain)
	resp := &response{
		context:          ctx,
		WriterChain:      chain,
		CleanPowerWriter: pw,
		header:           w.Header(),
	}
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp
}

// destroy destroy all reference that response keep
func (resp *response) destroy() {
	resp.context = nil
	resp.CleanPowerWriter = nil
	resp.WriterChain = nil
	resp.header = nil
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

// Redirect redirect to new url
func (resp *response) Redirect(url string) {
	http.Redirect(resp.httpWriter(), resp.request, url, http.StatusTemporaryRedirect)
}

// PermanentRedirect permanently redirect current request url to new url
func (resp *response) PermanentRedirect(url string) {
	http.Redirect(resp.httpWriter(), resp.request, url, http.StatusMovedPermanently)
}

// ReportStatus report an http status with given status code
func (resp *response) ReportStatus(statusCode int) {
	resp.BaseWriter().(http.ResponseWriter).WriteHeader(statusCode)
}

// Hijack hijack response connection
func (resp *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, is := resp.w.(http.Hijacker); is {
		return hijacker.Hijack()
	}
	return nil, nil, Err("Connection not support hijack")
}

// Flush flush response's output
func (resp *response) Flush() {
	if flusher, is := resp.w.(http.Flusher); is {
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
	resp.SetHeader(HEADER_SETCOOKIE, resp.newCookie(name, value, lifetime))
}

// SetSecureCookie setup response cookie with secureity
func (resp *response) SetSecureCookie(name, value string, lifetime int) {
	resp.SetCookie(name, resp.encodeSecureCookie(value), lifetime)
}

// DeleteClientCookie delete user briwser's cookie by name
func (resp *response) DeleteClientCookie(name string) {
	resp.SetCookie(name, "", -1)
}

// setSessionCookie setup session cookie, if enabled secure cookie, it will use it
func (resp *response) setSessionCookie(id string) {
	resp.SetSecureCookie(_COOKIE_SESSION, id, 0)
}

// Render render template with context
func (resp *response) Render(tmpl string) error {
	return resp.Server().RenderTemplate(resp, tmpl, resp.context)
}

// RenderWith render template with given value
func (resp *response) RenderWith(tmpl string, value interface{}) error {
	return resp.Server().RenderTemplate(resp, tmpl, value)
}

func (resp *response) httpWriter() http.ResponseWriter {
	return resp.BaseWriter().(http.ResponseWriter)
}
