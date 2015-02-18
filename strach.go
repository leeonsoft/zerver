// this file is strach before server start, and will be destroied
// at the moment of server init
package zerver_rest

type _strach struct {
	_funcHandlers map[string]*funcHandler // function handlers
}

var strach = &_strach{
	_funcHandlers: make(map[string]*funcHandler),
}

func (s *_strach) destroy() {
	s._funcHandlers = nil
}

func (s *_strach) funcHandler(pattern string) *funcHandler {
	return s._funcHandlers[pattern]
}

func (s *_strach) setFuncHandler(pattern string, handler *funcHandler) {
	s._funcHandlers[pattern] = handler
}
