package zerver

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/cosiner/golib/test"

	. "github.com/cosiner/golib/errors"
)

func TestCompile(t *testing.T) {
	tt := test.WrapTest(t)
	tt.Log(len(strings.Split("", "/")))
	tt.Log(compile("/"))
	tt.Log(compile("/:user/:id/:a"))
	tt.Log(compile("/user/:a/:abc/"))
	tt.Log(compile("/user/:/:abc/"))
}

// routes is copy from github.com/julienschmidt/go-http-routing-benchmark
func rt() *router {
	processor := new(routeProcessor)
	node := &router{str: "/user/id", processor: processor}
	fn := func(n *router, success bool) {
		if success {
			n.processor = processor
		}
	}
	fn(node.addPath("/user/i"))
	fn(node.addPath("/user/ie"))
	fn(node.addPath("/user/ief"))
	fn(node.addPath("/user/ieg"))
	fn(node.addPath("/title/|"))
	fn(node.addPath("/title/id/|"))
	fn(node.addPath("/title/i/|"))
	fn(node.addPath("/title/id/12"))
	fn(node.addPath("/ti/id/12"))
	fn(node.addPath("/ti/|/12"))

	// OAuth Authorizations
	fn(node.addPath("/authorizations"))
	fn(node.addPath("/authorizations/|"))
	fn(node.addPath("/authorizations"))

	fn(node.addPath("/authorizations/|"))
	fn(node.addPath("/applications/|/tokens/|"))
	fn(node.addPath("/applications/|/tokens"))
	fn(node.addPath("/applications/|/tokens/|"))
	// Activity
	fn(node.addPath("/events"))
	fn(node.addPath("/repos/|/|/events"))
	fn(node.addPath("/networks/|/|/events"))
	fn(node.addPath("/orgs/|/events"))
	fn(node.addPath("/users/|/received_events"))
	fn(node.addPath("/users/|/received_events/public"))
	fn(node.addPath("/users/|/events"))
	fn(node.addPath("/users/|/events/public"))
	fn(node.addPath("/users/|/events/orgs/|"))
	fn(node.addPath("/feeds"))
	fn(node.addPath("/notifications"))
	fn(node.addPath("/repos/|/|/notifications"))
	fn(node.addPath("/notifications"))
	fn(node.addPath("/repos/|/|/notifications"))
	fn(node.addPath("/notifications/threads/|"))

	fn(node.addPath("/notifications/threads/|/subscription"))
	fn(node.addPath("/notifications/threads/|/subscription"))
	fn(node.addPath("/notifications/threads/|/subscription"))
	fn(node.addPath("/repos/|/|/stargazers"))
	fn(node.addPath("/users/|/starred"))
	fn(node.addPath("/user/starred"))
	fn(node.addPath("/user/starred/|/|"))
	fn(node.addPath("/user/starred/|/|"))
	fn(node.addPath("/user/starred/|/|"))
	fn(node.addPath("/repos/|/|/subscribers"))
	fn(node.addPath("/users/|/subscriptions"))
	fn(node.addPath("/user/subscriptions"))
	fn(node.addPath("/repos/|/|/subscription"))
	fn(node.addPath("/repos/|/|/subscription"))
	fn(node.addPath("/repos/|/|/subscription"))
	fn(node.addPath("/user/subscriptions/|/|"))
	fn(node.addPath("/user/subscriptions/|/|"))
	fn(node.addPath("/user/subscriptions/|/|"))
	// Gists
	fn(node.addPath("/users/|/gists"))
	fn(node.addPath("/gists"))
	fn(node.addPath("/gists/public"))
	fn(node.addPath("/gists/starred"))
	fn(node.addPath("/gists/|"))
	fn(node.addPath("/gists"))

	fn(node.addPath("/gists/|/star"))
	fn(node.addPath("/gists/|/star"))
	fn(node.addPath("/gists/|/star"))
	fn(node.addPath("/gists/|/forks"))
	fn(node.addPath("/gists/|"))
	// Git Data
	fn(node.addPath("/repos/|/|/git/blobs/|"))
	fn(node.addPath("/repos/|/|/git/blobs"))
	fn(node.addPath("/repos/|/|/git/commits/|"))
	fn(node.addPath("/repos/|/|/git/commits"))
	fn(node.addPath("/repos/|/|/git/refs/|ref"))
	fn(node.addPath("/repos/|/|/git/refs"))
	fn(node.addPath("/repos/|/|/git/refs"))

	fn(node.addPath("/repos/|/|/git/tags/|"))
	fn(node.addPath("/repos/|/|/git/tags"))
	fn(node.addPath("/repos/|/|/git/trees/|"))
	fn(node.addPath("/repos/|/|/git/trees"))
	// Issues
	fn(node.addPath("/issues"))
	fn(node.addPath("/user/issues"))
	fn(node.addPath("/orgs/|/issues"))
	fn(node.addPath("/repos/|/|/issues"))
	fn(node.addPath("/repos/|/|/issues/|"))
	fn(node.addPath("/repos/|/|/issues"))

	fn(node.addPath("/repos/|/|/assignees"))
	fn(node.addPath("/repos/|/|/assignees/|"))
	fn(node.addPath("/repos/|/|/issues/|/comments"))
	fn(node.addPath("/repos/|/|/issues/comments"))
	fn(node.addPath("/repos/|/|/issues/comments/|"))
	fn(node.addPath("/repos/|/|/issues/|/comments"))

	fn(node.addPath("/repos/|/|/issues/|/events"))
	fn(node.addPath("/repos/|/|/issues/events"))
	fn(node.addPath("/repos/|/|/issues/events/|"))
	fn(node.addPath("/repos/|/|/labels"))
	fn(node.addPath("/repos/|/|/labels/|"))
	fn(node.addPath("/repos/|/|/labels"))

	fn(node.addPath("/repos/|/|/labels/|"))
	fn(node.addPath("/repos/|/|/issues/|/labels"))
	fn(node.addPath("/repos/|/|/issues/|/labels"))
	fn(node.addPath("/repos/|/|/issues/|/labels/|"))
	fn(node.addPath("/repos/|/|/issues/|/labels"))
	fn(node.addPath("/repos/|/|/issues/|/labels"))
	fn(node.addPath("/repos/|/|/milestones/|/labels"))
	fn(node.addPath("/repos/|/|/milestones"))
	fn(node.addPath("/repos/|/|/milestones/|"))
	fn(node.addPath("/repos/|/|/milestones"))

	fn(node.addPath("/repos/|/|/milestones/|"))
	// Miscellaneous
	fn(node.addPath("/emojis"))
	fn(node.addPath("/gitignore/templates"))
	fn(node.addPath("/gitignore/templates/|"))
	fn(node.addPath("/markdown"))
	fn(node.addPath("/markdown/raw"))
	fn(node.addPath("/meta"))
	fn(node.addPath("/rate_limit"))
	// Organizations
	fn(node.addPath("/users/|/orgs"))
	fn(node.addPath("/user/orgs"))
	fn(node.addPath("/orgs/|"))

	fn(node.addPath("/orgs/|/members"))
	fn(node.addPath("/orgs/|/members/|"))
	fn(node.addPath("/orgs/|/members/|"))
	fn(node.addPath("/orgs/|/public_members"))
	fn(node.addPath("/orgs/|/public_members/|"))
	fn(node.addPath("/orgs/|/public_members/|"))
	fn(node.addPath("/orgs/|/public_members/|"))
	fn(node.addPath("/orgs/|/teams"))
	fn(node.addPath("/teams/|"))
	fn(node.addPath("/orgs/|/teams"))

	fn(node.addPath("/teams/|"))
	fn(node.addPath("/teams/|/members"))
	fn(node.addPath("/teams/|/members/|"))
	fn(node.addPath("/teams/|/members/|"))
	fn(node.addPath("/teams/|/members/|"))
	fn(node.addPath("/teams/|/repos"))
	fn(node.addPath("/teams/|/repos/|/|"))
	fn(node.addPath("/teams/|/repos/|/|"))
	fn(node.addPath("/teams/|/repos/|/|"))
	fn(node.addPath("/user/teams"))
	// Pull Requests
	fn(node.addPath("/repos/|/|/pulls"))
	fn(node.addPath("/repos/|/|/pulls/|"))
	fn(node.addPath("/repos/|/|/pulls"))

	fn(node.addPath("/repos/|/|/pulls/|/commits"))
	fn(node.addPath("/repos/|/|/pulls/|/files"))
	fn(node.addPath("/repos/|/|/pulls/|/merge"))
	fn(node.addPath("/repos/|/|/pulls/|/merge"))
	fn(node.addPath("/repos/|/|/pulls/|/comments"))
	fn(node.addPath("/repos/|/|/pulls/comments"))
	fn(node.addPath("/repos/|/|/pulls/comments/|"))
	fn(node.addPath("/repos/|/|/pulls/|/comments"))

	// Repositories
	fn(node.addPath("/user/repos"))
	fn(node.addPath("/users/|/repos"))
	fn(node.addPath("/orgs/|/repos"))
	fn(node.addPath("/repositories"))
	fn(node.addPath("/user/repos"))
	fn(node.addPath("/orgs/|/repos"))
	fn(node.addPath("/repos/|/|"))

	fn(node.addPath("/repos/|/|/contributors"))
	fn(node.addPath("/repos/|/|/languages"))
	fn(node.addPath("/repos/|/|/teams"))
	fn(node.addPath("/repos/|/|/tags"))
	fn(node.addPath("/repos/|/|/branches"))
	fn(node.addPath("/repos/|/|/branches/|"))
	fn(node.addPath("/repos/|/|"))
	fn(node.addPath("/repos/|/|/collaborators"))
	fn(node.addPath("/repos/|/|/collaborators/|"))
	fn(node.addPath("/repos/|/|/collaborators/|"))
	fn(node.addPath("/repos/|/|/collaborators/|"))
	fn(node.addPath("/repos/|/|/comments"))
	fn(node.addPath("/repos/|/|/commits/|/comments"))
	fn(node.addPath("/repos/|/|/commits/|/comments"))
	fn(node.addPath("/repos/|/|/comments/|"))

	fn(node.addPath("/repos/|/|/comments/|"))
	fn(node.addPath("/repos/|/|/commits"))
	fn(node.addPath("/repos/|/|/commits/|"))
	fn(node.addPath("/repos/|/|/readme"))
	fn(node.addPath("/repos/|/|/contents/|path"))

	fn(node.addPath("/repos/|/|/|/|"))
	fn(node.addPath("/repos/|/|/keys"))
	fn(node.addPath("/repos/|/|/keys/|"))
	fn(node.addPath("/repos/|/|/keys"))

	fn(node.addPath("/repos/|/|/keys/|"))
	fn(node.addPath("/repos/|/|/downloads"))
	fn(node.addPath("/repos/|/|/downloads/|"))
	fn(node.addPath("/repos/|/|/downloads/|"))
	fn(node.addPath("/repos/|/|/forks"))
	fn(node.addPath("/repos/|/|/forks"))
	fn(node.addPath("/repos/|/|/hooks"))
	fn(node.addPath("/repos/|/|/hooks/|"))
	fn(node.addPath("/repos/|/|/hooks"))

	fn(node.addPath("/repos/|/|/hooks/|/tests"))
	fn(node.addPath("/repos/|/|/hooks/|"))
	fn(node.addPath("/repos/|/|/merges"))
	fn(node.addPath("/repos/|/|/releases"))
	fn(node.addPath("/repos/|/|/releases/|"))
	fn(node.addPath("/repos/|/|/releases"))

	fn(node.addPath("/repos/|/|/releases/|"))
	fn(node.addPath("/repos/|/|/releases/|/assets"))
	fn(node.addPath("/repos/|/|/stats/contributors"))
	fn(node.addPath("/repos/|/|/stats/commit_activity"))
	fn(node.addPath("/repos/|/|/stats/code_frequency"))
	fn(node.addPath("/repos/|/|/stats/participation"))
	fn(node.addPath("/repos/|/|/stats/punch_card"))
	fn(node.addPath("/repos/|/|/statuses/|"))
	fn(node.addPath("/repos/|/|/statuses/|"))
	// Search
	fn(node.addPath("/search/repositories"))
	fn(node.addPath("/search/code"))
	fn(node.addPath("/search/issues"))
	fn(node.addPath("/search/users"))
	fn(node.addPath("/legacy/issues/search/|/|/|/|"))
	fn(node.addPath("/legacy/repos/search/|"))
	fn(node.addPath("/legacy/user/search/|"))
	fn(node.addPath("/legacy/user/email/|"))
	// Users
	fn(node.addPath("/users/|"))
	fn(node.addPath("/user"))

	fn(node.addPath("/users"))
	fn(node.addPath("/user/emails"))
	fn(node.addPath("/user/emails"))
	fn(node.addPath("/user/emails"))
	fn(node.addPath("/users/|/followers"))
	fn(node.addPath("/user/followers"))
	fn(node.addPath("/users/|/following"))
	fn(node.addPath("/user/following"))
	fn(node.addPath("/user/following/|"))
	fn(node.addPath("/users/|/following/|"))
	fn(node.addPath("/user/following/|"))
	fn(node.addPath("/user/following/|"))
	fn(node.addPath("/users/|/keys"))
	fn(node.addPath("/user/keys"))
	fn(node.addPath("/user/keys/|"))
	fn(node.addPath("/user/keys"))

	fn(node.addPath("/user/keys/|"))
	// PrintRouteTree(node)
	return node
}

var r = rt()

func BenchmarkMatchRouteOne(b *testing.B) {
	// tt := test.WrapTest(b)
	// path := "/repos/julienschmidt/httprouter/stargazers"
	path := "/user/repos"
	// path := "/user/keys"
	// path := "/user/aa/exist"
	for i := 0; i < b.N; i++ {
		// pathIndex := 0
		// var vars []string
		// var continu = true
		// n := r
		// for continu {
		// 	pathIndex, vars, n, continu = n.matchMulti(path, pathIndex, vars)
		// }
		_, _ = r.matchOne(path)
	}
}

func BenchmarkMakeSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make([]string, 0, PathVarCount)
	}
}

func BenchmarkMatchRouteMultiple(b *testing.B) {
	// tt := test.WrapTest(b)
	// path := "/legacy/issues/search/aaa/bbb/ccc/ddd"
	// path := "/user/repos"
	path := "/repos/julienschmidt/httprouter/stargazers"
	// path := "/user/aa/exist"
	for i := 0; i < b.N; i++ {
		pathIndex := 0
		var vars []string
		var continu = true
		n := r
		for continu {
			pathIndex, vars, n, continu = n.matchMultiple(path, pathIndex, vars)
		}
		if n == nil {
			b.Fail()
		}
	}
}

func TestRoute(t *testing.T) {
	rt := new(router)
	OnErrPanic(rt.AddHandler("/user.:format", newFuncHandler()))
	OnErrPanic(rt.AddHandler("/v:version", newFuncHandler()))
	OnErrPanic(rt.AddHandler("/vaa/:id", newFuncHandler()))
	// OnErrPanic(rt.AddHandler("/vba/:id", newFuncHandler()))
	// OnErrPanic(rt.AddHandler("/v0a/:id", newFuncHandler()))
	rt.PrintRouteTree(os.Stdout)
	_, value := rt.matchOne("/user.json")
	t.Log(value)
	_, value = rt.matchOne("/vbc")
	t.Log(value)
}

func TestFilterHideHandler(t *testing.T) {
	tt := test.WrapTest(t)
	rt := NewRouter()
	rt.AddFuncFilter("/user/12:id", EmptyFilterFunc)
	rt.AddFuncHandler("/user/:id", "GET", EmptyHandlerFunc)
	h, _, _ := rt.MatchHandlerFilters(&url.URL{Path: "/user/1234"})
	tt.AssertTrue(h == nil)
	h, _, _ = rt.MatchHandlerFilters(&url.URL{Path: "/user/2234"})
	tt.AssertTrue(h != nil)
}

func TestConflictWildcardCatchall(t *testing.T) {
	tt := test.WrapTest(t)
	rt := new(router)
	tt.AssertNil(rt.AddFuncHandler("/:user/:id", "GET", EmptyHandlerFunc))
	e := rt.AddFuncHandler("/*user", "GET", EmptyHandlerFunc)
	tt.AssertTrue(e != nil)
}
