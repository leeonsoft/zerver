package filters

import (
	"compress/flate"
	"compress/gzip"
	"strings"

	"github.com/cosiner/zerver"
)

type gzipResponse struct {
	gzipWriter *gzip.Writer
	zerver.Response
}

type flateResponse struct {
	flateWriter *flate.Writer
	zerver.Response
}

func (gr *gzipResponse) Write(data []byte) (int, error) {
	return gr.gzipWriter.Write(data)
}

func (fr *flateResponse) Write(data []byte) (int, error) {
	return fr.flateWriter.Write(data)
}

func NewCompressFilter() zerver.Filter {
	return (zerver.FilterFunc)(compressFilter)
}

func compressFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	encoding := req.AcceptEncodings()
	if strings.Contains(encoding, zerver.ENCODING_GZIP) {
		gzw := gzip.NewWriter(resp)
		resp.SetContentEncoding(zerver.ENCODING_GZIP)
		chain(req, &gzipResponse{gzw, resp})
		gzw.Close()
	} else if strings.Contains(encoding, zerver.ENCODING_DEFLATE) {
		flw, _ := flate.NewWriter(resp, flate.DefaultCompression)
		resp.SetContentEncoding(zerver.ENCODING_DEFLATE)
		chain(req, &flateResponse{flw, resp})
		flw.Close()
	} else {
		chain(req, resp)
	}
}
