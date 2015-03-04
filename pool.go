package zerver_rest

import "sync"

type serverPool struct {
	requestPool    *sync.Pool
	responsePool   *sync.Pool
	varIndexerPool *sync.Pool
	varsPool       *sync.Pool
}

var pool *serverPool = &serverPool{
	requestPool: &sync.Pool{
		New: func() interface{} {
			return new(request)
		},
	},
	responsePool: &sync.Pool{
		New: func() interface{} {
			return new(response)
		},
	},
	varIndexerPool: &sync.Pool{
		New: func() interface{} {
			return new(urlVarIndexer)
		},
	},
	varsPool: &sync.Pool{
		New: func() interface{} {
			return make([]string, 0, PathVarCount)
		},
	},
}

func (pool *serverPool) newRequest() *request {
	return pool.requestPool.Get().(*request)
}

func (pool *serverPool) newResponse() *response {
	return pool.responsePool.Get().(*response)
}

func (pool *serverPool) newVarIndexer() *urlVarIndexer {
	return pool.varIndexerPool.Get().(*urlVarIndexer)
}

func (pool *serverPool) newVars() []string {
	return pool.varsPool.Get().([]string)
}

func (pool *serverPool) recycleRequest(req *request) {
	pool.requestPool.Put(req)
}

func (pool *serverPool) recycleResponse(resp *response) {
	pool.responsePool.Put(resp)
}

func (pool *serverPool) recycleVarIndexer(indexer *urlVarIndexer) {
	pool.varIndexerPool.Put(indexer)
}

func (pool *serverPool) recycleVars(vars []string) {
	pool.varsPool.Put(vars[:0])
}
