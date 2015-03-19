package filters

import (
	"bytes"
	"encoding/base64"
	"strings"

	"github.com/cosiner/zerver"
)

type BasicAuthFilter struct {
	AuthUserAttrName string
	AuthPassAttrName string
}

func (b *BasicAuthFilter) Init(*zerver.Server) error {
	if b.AuthUserAttrName == "" {
		b.AuthUserAttrName = "AuthUser"
	}
	if b.AuthPassAttrName == "" {
		b.AuthPassAttrName = "AuthPass"
	}
	return nil
}

func (b *BasicAuthFilter) Filter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	auth := strings.TrimSpace(req.Header(_HEADER_AUTHRIZATION))
	if auth == "" {
		resp.ReportUnauthorized()
		return
	}
	authStr, err := base64.URLEncoding.DecodeString(auth)
	if err != nil {
		resp.ReportBadRequest()
		return
	}
	index := bytes.IndexByte(authStr, ':')
	if index <= 0 || index == len(authStr)-1 {
		resp.ReportBadRequest()
		return
	}
	req.SetAttr(b.AuthUserAttrName, authStr[:index])
	req.SetAttr(b.AuthPassAttrName, authStr[index+1:])
	chain(req, resp)
}

func (b *BasicAuthFilter) Destroy() {}
