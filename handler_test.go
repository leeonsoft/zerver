package zerver_rest

var methods = []string{GET, POST, PUT, DELETE, PATCH}

func fhandler() Handler {
	processFn := func(Request, Response) {}
	fh := newFuncHandler()
	for i := 0; i < len(methods); i++ {
		fh.setMethodHandler(methods[i], processFn)
	}
	return fh
}
