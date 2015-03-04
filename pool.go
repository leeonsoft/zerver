package zerver_rest

import (
	. "github.com/cosiner/golib/errors"

	"sync"
)

type ServerPool struct {
	requestPool    *sync.Pool
	responsePool   *sync.Pool
	varIndexerPool *sync.Pool
	varsPool       *sync.Pool
	otherPools     map[string]*sync.Pool
}

var Pool *ServerPool = &ServerPool{
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
	otherPools: make(map[string]*sync.Pool),
}

func (pool *ServerPool) ReigisterPool(name string, newFunc func() interface{}) error {
	op := pool.otherPools
	if _, has := op[name]; has {
		return Err("Pool for " + name + " already exist")
	}
	op[name] = &sync.Pool{New: newFunc}
}

func (pool *ServerPool) NewFrom(name string) interface{} {
	return pool.otherPools[name].Get()
}

func (pool *ServerPool) newRequest() *request {
	return pool.requestPool.Get().(*request)
}

func (pool *ServerPool) newResponse() *response {
	return pool.responsePool.Get().(*response)
}

func (pool *ServerPool) newVarIndexer() *urlVarIndexer {
	return pool.varIndexerPool.Get().(*urlVarIndexer)
}

func (pool *ServerPool) newVars() []string {
	return pool.varsPool.Get().([]string)
}

func (pool *ServerPool) recycleRequest(req *request) {
	pool.requestPool.Put(req)
}

func (pool *ServerPool) recycleResponse(resp *response) {
	pool.responsePool.Put(resp)
}

func (pool *ServerPool) recycleVarIndexer(indexer *urlVarIndexer) {
	pool.varIndexerPool.Put(indexer)
}

func (pool *ServerPool) recycleVars(vars []string) {
	pool.varsPool.Put(vars[:0])
}

func (pool *ServerPool) RecycleTo(name string, value interface{}) {
	pool.otherPools[name].Put(value)
}
