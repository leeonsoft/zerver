package router

import (
	"net/url"
	"os"
	"testing"

	"github.com/cosiner/zerver"

	"github.com/cosiner/golib/test"
)

var rt = func() zerver.Router {
	r := zerver.NewRouter()
	gr := NewGroupRouter("/user", r)
	gr.Get("/:id", zerver.EmptyHandlerFunc)
	gr.Post("/:id", zerver.EmptyHandlerFunc)
	gr.Delete("/:id", zerver.EmptyHandlerFunc)
	gr.Post("/:id/exist", zerver.EmptyHandlerFunc)
	return r
}()
var u = &url.URL{Path: "/user/123/exist"}

func TestGroup(t *testing.T) {
	tt := test.WrapTest(t)
	rt.PrintRouteTree(os.Stdout)
	h, _, _ := rt.MatchHandlerFilters(u)
	tt.AssertTrue(h != nil)
}

func BenchmarkRouter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = rt.MatchHandlerFilters(u)
	}
}
