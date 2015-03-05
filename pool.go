package zerver_rest

import (
	. "github.com/cosiner/golib/errors"

	"sync"
)

type requestEnv struct {
	req  *request
	resp *response
}

type ServerPool struct {
	requestEnvPool *sync.Pool
	varIndexerPool *sync.Pool
	filtersPool    *sync.Pool
	otherPools     map[string]*sync.Pool
}

var Pool *ServerPool = &ServerPool{
	requestEnvPool: &sync.Pool{
		New: func() interface{} {
			return &requestEnv{new(request), new(response)}
		},
	},
	varIndexerPool: &sync.Pool{
		New: func() interface{} {
			return &urlVarIndexer{values: make([]string, 0, PathVarCount)}
		},
	},
	filtersPool: &sync.Pool{
		New: func() interface{} {
			return make([]Filter, 0, FilterCount)
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
	return nil
}

func (pool *ServerPool) NewFrom(name string) interface{} {
	return pool.otherPools[name].Get()
}

func (pool *ServerPool) newRequestEnv() *requestEnv {
	return pool.requestEnvPool.Get().(*requestEnv)
}

func (pool *ServerPool) newVarIndexer() *urlVarIndexer {
	return pool.varIndexerPool.Get().(*urlVarIndexer)
}

func (pool *ServerPool) newFilters() []Filter {
	return pool.filtersPool.Get().([]Filter)
}

func (pool *ServerPool) recycleRequestEnv(req *requestEnv) {
	pool.requestEnvPool.Put(req)
}

func (pool *ServerPool) recycleVarIndexer(indexer URLVarIndexer) {
	pool.varIndexerPool.Put(indexer)
}

func (pool *ServerPool) recycleFilters(filters []Filter) {
	if filters != nil {
		filters = filters[:0]
		pool.filtersPool.Put(filters)
	}
}

func (pool *ServerPool) RecycleTo(name string, value interface{}) {
	pool.otherPools[name].Put(value)
}
