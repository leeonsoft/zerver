package pprof

import (
	"testing"

	"github.com/cosiner/zerver"
)

func TestMonitor(t *testing.T) {
	rt := zerver.NewRouter()
	filters := zerver.NewRootFilters()
	EnableMonitoring("/pprof", rt, filters)
	server := zerver.NewServerWith(rt, filters)
	server.Start(":4000")
}
