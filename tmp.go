package zerver_rest

type _tmp struct {
	_funcHandlers map[string]*funcHandler // function handlers - pattern->handler
}

var tmp = &_tmp{
	_funcHandlers: make(map[string]*funcHandler),
}

func (t *_tmp) destroy() {
	t._funcHandlers = nil
}

func (t *_tmp) funcHandler(pattern string) *funcHandler {
	return t._funcHandlers[pattern]
}

func (t *_tmp) setFuncHandler(pattern string, handler *funcHandler) {
	t._funcHandlers[pattern] = handler
}
