package scheduler

// Scheduler interface
type Scheduler interface {
	AddFunc(spec string, cmd func()) error
	Start() error
}
