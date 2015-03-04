// this file is tmp before server start, and will be destroied
// at the moment of server init
package zerver_rest

type _tmp struct {
	_funcHandlers map[string]*funcHandler // function handlers
}

var tmp = &_strach{
	_funcHandlers: make(map[string]*funcHandler),
}

func (s *_tmp) destroy() {
	s._funcHandlers = nil
}

func (s *_tmp) funcHandler(pattern string) *funcHandler {
	return s._funcHandlers[pattern]
}

func (s *_tmp) setFuncHandler(pattern string, handler *funcHandler) {
	s._funcHandlers[pattern] = handler
}
