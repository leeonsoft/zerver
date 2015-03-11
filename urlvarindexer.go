package zerver

type (
	// URLVarIndexer is a indexer for name to value
	URLVarIndexer interface {
		// URLVar return value of variable
		URLVar(name string) string
		// ScanURLVars scan given values into variable addresses
		// if address is nil, skip it
		ScanURLVars(vars ...*string)
		// URLVars return all values of variable
		URLVars() []string
		destroySelf() // avoid confilict with Request interface
	}

	// urlVarIndexer is an implementation of URLVarIndexer
	urlVarIndexer struct {
		vars   map[string]int // url variables and indexs of sections splited by '/'
		values []string       // all url variable values
	}
)

func (v *urlVarIndexer) destroySelf() {
	v.values = v.values[:0]
	v.vars = nil
	Pool.recycleVarIndexer(v)
}

// URLVar return values of variable
func (v *urlVarIndexer) URLVar(name string) string {
	if len(v.vars) > 0 {
		if index, has := v.vars[name]; has {
			return v.values[index]
		}
	}
	return ""
}

// URLVars return all values of variable
func (v *urlVarIndexer) URLVars() []string {
	if len(v.vars) > 0 {
		return v.values
	}
	return nil
}

// ScanURLVars scan values into variable addresses
// if address is nil, skip it
func (v *urlVarIndexer) ScanURLVars(vars ...*string) {
	if len(v.vars) > 0 {
		values := v.values
		l1, l2 := len(values), len(vars)
		for i := 0; i < l1 && i < l2; i++ {
			if vars[i] != nil {
				*vars[i] = values[i]
			}
		}
	}
}
