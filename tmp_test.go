package zerver

import "testing"

func TestTmp(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fail()
		}
	}()
	tmpDestroy()
	TmpSet("test", "aa")
}
