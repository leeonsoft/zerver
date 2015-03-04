package zerver_rest

type (
	// Task
	Task interface {
		URLVarIndexer
		Value() interface{}
	}

	task struct {
		URLVarIndexer
		value interface{}
	}

	TaskHandlerFunc func(Task)

	TaskHandler interface {
		Init(*Server) error
		Destroy()
		Handle(Task)
	}
)

func (t *task) Value() interface{} {
	return t.value
}

func newTask(indexer URLVarIndexer, value interface{}) Task {
	return &task{
		URLVarIndexer: indexer,
		value:         value,
	}
}

func (TaskHandlerFunc) Init(*Server) error  { return nil }
func (fn TaskHandlerFunc) Handle(task Task) { fn(task) }
func (TaskHandlerFunc) Destroy()            {}
