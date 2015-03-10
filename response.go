package zerver_rest

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	. "github.com/cosiner/golib/errors"
)

type (
	Response interface {
		SetHeader(name, value string)
		AddHeader(name, value string)
		RemoveHeader(name string)
		SetContentEncoding(enc string)
		SetContentType(typ string)
		SetAdvancedCookie(c *http.Cookie)
		SetCookie(name, value string, lifetime int)
		SetSecureCookie(name, value string, lifetime int)
		DeleteClientCookie(name string)
		ReportStatus(statusCode int)
		Hijack() (net.Conn, *bufio.ReadWriter, error)
		Flush()
		io.Writer
		StatusResponse
		CacheSeconds(secs int)
		CacheUntil(*time.Time)
		NoCache()
		destroy()
	}

	// response represent a response of request to user
	response struct {
		http.ResponseWriter
		header http.Header
	}
)

// newResponse create a new response, and set default content type to HTML
func (resp *response) init(w http.ResponseWriter) Response {
	resp.ResponseWriter = w
	resp.header = w.Header()
	return resp
}

func (resp *response) destroy() {
	resp.ResponseWriter = nil
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

func (resp *response) CacheSeconds(secs int) {
	resp.SetHeader(HEADER_CACHECONTROL, fmt.Sprintf("max-age:%d", secs))
}

func (resp *response) CacheUntil(t *time.Time) {
	resp.SetHeader(HEADER_EXPIRES, t.Format(http.TimeFormat))
}

func (resp *response) NoCache() {
	resp.SetHeader(HEADER_CACHECONTROL, "no-cache")
}

func (resp *response) SetAdvancedCookie(c *http.Cookie) {
	resp.AddHeader(HEADER_SETCOOKIE, c.String())
}

// SetCookie setup response cookie
func (resp *response) SetCookie(name, value string, lifetime int) {
	resp.SetAdvancedCookie(&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: lifetime,
	})
}

// SetSecureCookie setup response cookie with secureity
func (resp *response) SetSecureCookie(name, value string, lifetime int) {
	resp.SetCookie(name, value, lifetime)
}

// DeleteClientCookie delete user briwser's cookie by name
func (resp *response) DeleteClientCookie(name string) {
	resp.SetCookie(name, "", -1)
}
