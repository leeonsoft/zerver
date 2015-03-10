package zerver_rest

import "testing"

var methods = []string{GET, POST, PUT, DELETE, PATCH}

func fhandler() Handler {
	processFn := func(Request, Response) {}
	fh := newFuncHandler()
	for i := 0; i < len(methods); i++ {
		fh.setMethodHandler(methods[i], processFn)
	}
	return fh
}

func BenchmarkIndicateFuncHandler(b *testing.B) {
	index := 0
	fh := fhandler()
	for i := 0; i < b.N; i++ {
		_ = IndicateHandler(methods[index], fh)
		index++
		if index == len(methods) {
			index = 0
		}
	}
}

func BenchmarkIndicateHandler(b *testing.B) {
	index := 0
	fh := &EmptyHandler{}
	for i := 0; i < b.N; i++ {
		_ = standardIndicate(methods[index], fh)
		index++
		if index == len(methods) {
			index = 0
		}
	}
}
