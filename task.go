package zerver_rest

type TaskHandlerFunc func(Request)
type TaskHandler interface {
	Init(*Server) error
	Destroy()
	Handle(Request)
}

func (TaskHandlerFunc) Init(*Server) error    { return nil }
func (fn TaskHandlerFunc) Handle(req Request) { fn(req) }
func (TaskHandlerFunc) Destroy()              {}
