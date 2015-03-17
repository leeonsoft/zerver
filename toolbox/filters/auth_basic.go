package filters

import (
	"bytes"
	"encoding/base64"
	"strings"

	"github.com/cosiner/zerver"
)

var (
	AuthUserAttrName string
	AuthPassAttrName string
)

func BasicAuthFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	auth := strings.TrimSpace(req.Header(_HEADER_AUTHRIZATION))
	if auth == "" {
		resp.ReportUnauthorized()
		return
	}
	authStr, err := base64.URLEncoding.DecodeString(auth)
	if err != nil {
		resp.ReportUnauthorized()
		return
	}
	index := bytes.IndexByte(authStr, ':')
	if index <= 0 || index == len(authStr)-1 {
		resp.ReportUnauthorized()
		return
	}
	req.SetAttr(AuthUserAttrName, authStr[:index])
	req.SetAttr(AuthPassAttrName, authStr[index+1:])
	chain(req, resp)
}
