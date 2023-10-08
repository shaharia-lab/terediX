// Package scheduler helps to schedule the task
package scheduler

// Scheduler interface
type Scheduler interface {
	AddFunc(spec string, cmd func()) error
	Start() error
}
