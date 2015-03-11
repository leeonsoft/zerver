package pprof

import "testing"

func TestMonitor(t *testing.T) {
	NewMonitorServer("/").Start(":4000")
}
