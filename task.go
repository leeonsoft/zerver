package zerver

type (
	// Task
	Task interface {
		URLVarIndexer
		Value() interface{}
		serverGetter
	}

	TaskHandlerFunc func(Task)

	TaskHandler interface {
		Init(*Server) error
		Destroy()
		Handle(Task)
	}

	task struct {
		serverGetter
		URLVarIndexer
		value interface{}
	}
)

func newTask(s serverGetter, indexer URLVarIndexer, value interface{}) Task {
	return &task{
		serverGetter:  s,
		URLVarIndexer: indexer,
		value:         value,
	}
}

func (t *task) Value() interface{} {
	return t.value
}

func (TaskHandlerFunc) Init(*Server) error  { return nil }
func (fn TaskHandlerFunc) Handle(task Task) { fn(task) }
func (TaskHandlerFunc) Destroy()            {}
