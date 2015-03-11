package filters

import (
	"fmt"

	"github.com/cosiner/golib/types"
	"github.com/cosiner/zerver"
)

func ErrorToBrowserFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	defer func() {
		if e := recover(); e != nil {
			resp.Write(types.UnsafeBytes(fmt.Sprint(e)))
		}
	}()
	chain(req, resp)
}
