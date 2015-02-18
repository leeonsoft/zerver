package zerver_rest

import "strings"

const (
	// Http Header
	HEADER_CONTENTTYPE     = "Content-Type"
	HEADER_CONTENTLENGTH   = "Content-Length"
	HEADER_SETCOOKIE       = "Set-Cookie"
	HEADER_REFER           = "Referer"
	HEADER_CONTENTENCODING = "Content-Encoding"
	HEADER_USERAGENT       = "User-Agent"
	HEADER_ACCEPTENCODING  = "Accept-Encoding"

	// ContentEncoding
	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"

	// Request Method
	GET            = "GET"
	POST           = "POST"
	DELETE         = "DELETE"
	PUT            = "PUT"
	PATCH          = "PATCH"
	UNKNOWN_METHOD = "UNKNOWN"

	// Content Type
	CONTNTTYPE_PLAIN = "text/plain"
	CONTENTTYPE_HTML = "text/html"
	CONTENTTYPE_XML  = "application/xml"
	CONTENTTYPE_JSON = "application/json"
)

// parseRequestMethod convert a string to request method, default use GET
// if string is empty
func parseRequestMethod(s string) string {
	if s == "" {
		return GET
	}
	return strings.ToUpper(s)
}

// parseContentType parse content type
func parseContentType(str string) string {
	if str == "" {
		return CONTENTTYPE_HTML
	}
	return strings.ToLower(strings.TrimSpace(str))
}
