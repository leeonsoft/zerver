package zerver

import (
	"compress/gzip"
	"io"
	"strings"
)

type (
	// ResponseCompressor enable compress for response
	ResponseCompressor interface {
		// EnableCompress enable compress for response, return compress writer
		// it should be closed after handle response by CloseResponseCompress
		EnableCompress(Request, Response) io.WriteCloser
		// CloseResponseCompress close compress writer
		CloseResponseCompress(Response, io.WriteCloser)
	}
	// emptyCompressor is an empty response compressor
	emptyCompressor struct{}
	// gzipCompressor enable gzip compress for response
	gzipCompressor struct{}
)

// NewResponseCompressor create a new response compressor, if not enabled,
// return an empty compressor
func NewResponseCompressor(enable bool) ResponseCompressor {
	if enable {
		return gzipCompressor{}
	}
	return emptyCompressor{}
}

func (ec emptyCompressor) EnableCompress(_ Request, _ Response) io.WriteCloser { return nil }
func (ec emptyCompressor) CloseResponseCompress(_ Response, _ io.WriteCloser)  {}

func (gc gzipCompressor) EnableCompress(req Request, resp Response) io.WriteCloser {
	var gzw *gzip.Writer
	if strings.Contains(req.AcceptEncodings(), ENCODING_GZIP) {
		gzw = gzip.NewWriter(resp.Writer())
		resp.SetContentEncoding(ENCODING_GZIP)
		resp.Wrap(gzw)
	}
	return gzw
}

func (gc gzipCompressor) CloseResponseCompress(resp Response, wr io.WriteCloser) {
	if wr != nil {
		resp.Unwrap()
		wr.Close()
	}
}
