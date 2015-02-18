package filters

import (
	"compress/flate"
	"compress/gzip"
	"strings"

	zerver "github.com/cosiner/zerver_rest"
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

func CompressFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	encoding := req.AcceptEncodings()
	if strings.Contains(encoding, zerver.ENCODING_GZIP) {
		gzw := gzip.NewWriter(resp)
		resp.SetContentEncoding(zerver.ENCODING_GZIP)
		chain.Filter(req, &gzipResponse{gzw, resp})
		gzw.Close()
	} else if strings.Contains(encoding, zerver.ENCODING_DEFLATE) {
		flw, _ := flate.NewWriter(resp, flate.DefaultCompression)
		resp.SetContentEncoding(zerver.ENCODING_DEFLATE)
		chain.Filter(req, &flateResponse{flw, resp})
		flw.Close()
	} else {
		chain.Filter(req, resp)
	}
}
