package zerver

import "sync"

type (
	// ValuesPool is a pool for url variable values
	VaulesPool interface {
		Require() []string
		Recycle([]string)
	}

	pathValuesPool struct {
		pool [][]string
		*sync.RWMutex
	}
)

var (
	valuesPool = newValuesPool(1024)
)

func newValuesPool(size int) VaulesPool {
	p := new(pathValuesPool)
	p.pool = make([][]string, 0, size)
	p.RWMutex = new(sync.RWMutex)
	return p
}

func (p *pathValuesPool) Require() (v []string) {
	p.Lock()
	pool := p.pool
	len := len(pool)
	if len == 0 {
		v = make([]string, 0, PathVarCount)
	} else {
		v = pool[len-1]
		p.pool = pool[:len-1]
	}
	p.Unlock()
	return
}

func (p *pathValuesPool) Recycle(v []string) {
	if cap(v) != 0 {
		v = v[:0]
		p.Lock()
		pool := p.pool
		if len := len(pool); len != cap(pool) { // if pool is full, do nothing
			pool = pool[:len+1]
			pool[len] = v
			p.pool = pool
		}
		p.Unlock()
	}
}
