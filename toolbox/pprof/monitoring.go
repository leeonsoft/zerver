package pprof

import (
	"log"
	"net/http"
	"net/url"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/cosiner/zerver_rest/toolbox/filters"

	"github.com/cosiner/golib/sys"
	"github.com/cosiner/golib/types"
	zerver "github.com/cosiner/zerver_rest"
)

var path = "/status"

func EnableMonitoring(p string, rt zerver.Router) {
	if p != "" {
		path = p
	}
	initRoutes()
	for subpath, handler := range routes {
		rt.AddFuncHandler(path+subpath, "GET", handler)
		options = append(options, "GET "+path+subpath+": "+infos[subpath]+"\n")
	}
	rt.AddFilter(path, filters.NewLogFilter(log.Println))
	rt.AddFuncFilter(path, globalFilter)
}

var options = make([]string, 0, len(infos)+1)
var routes = make(map[string]zerver.HandlerFunc)
var infos = make(map[string]string)

func pprofLookupHandler(name string) zerver.HandlerFunc {
	return func(req zerver.Request, resp zerver.Response) {
		pprof.Lookup(name).WriteTo(resp, 2)
	}
}

func globalFilter(req zerver.Request, resp zerver.Response, chain zerver.FilterChain) {
	resp.SetContentType("text/plain")
	if resp.Status() == http.StatusNotFound {
		resp.SetHeader("Location", path+"/options?from="+url.QueryEscape(req.URL().Path))
		resp.ReportStatus(http.StatusMovedPermanently)
	} else if resp.Status() == http.StatusMethodNotAllowed {
		sys.WriteStrln(resp, "The pprof interface only support GET request")
	} else {
		chain(req, resp)
	}
}

func initRoutes() {
	infos["/goroutine"] = "Get goroutine info"
	routes["/goroutine"] = pprofLookupHandler("goroutine")
	infos["/heap"] = "Get heap info"
	routes["/heap"] = pprofLookupHandler("heap")
	infos["/thread"] = "Get thread create info"
	routes["/thread"] = pprofLookupHandler("threadcreate")
	infos["/block"] = "Get block info"
	routes["/block"] = pprofLookupHandler("block")

	infos["/cpu"] = "Get CPU info, default seconds is 30, use ?seconds= to reset"
	routes["/cpu"] = func(req zerver.Request, resp zerver.Response) {
		var t int
		if secs := req.Param("seconds"); secs != "" {
			var err error
			if t, err = strconv.Atoi(secs); err != nil {
				resp.ReportBadRequest()
				sys.WriteStrln(resp, secs+" is not a integer number")
				return
			}
		}
		if t <= 0 {
			t = 30
		}
		pprof.StartCPUProfile(resp)
		time.Sleep(time.Duration(t) * time.Second)
		pprof.StopCPUProfile()
	}

	infos["/memory"] = "Get memory info"
	routes["/memory"] = func(req zerver.Request, resp zerver.Response) {
		runtime.GC()
		pprof.WriteHeapProfile(resp)
	}

	infos["/routes"] = "Get all routes"
	routes["/routes"] = func(req zerver.Request, resp zerver.Response) {
		req.Server().PrintRouteTree(resp)
	}

	infos["/options"] = "Get all pprof options"
	routes["/options"] = func(req zerver.Request, resp zerver.Response) {
		if from := req.Param("from"); from != "" {
			resp.Write(types.UnsafeBytes("There is no this pprof option: " + from + "\n"))
		}
		for i := range options {
			resp.Write(types.UnsafeBytes(options[i]))
		}
	}
}
