// Package scheduler helps to schedule the task
package scheduler

// StaticScheduler definition
type StaticScheduler struct {
	jobs []func()
}

// NewStaticScheduler to create new static scheduler
func NewStaticScheduler() *StaticScheduler {
	return &StaticScheduler{}
}

// AddFunc to add new function to cron scheduler
func (ss *StaticScheduler) AddFunc(spec string, cmd func()) error {
	ss.jobs = append(ss.jobs, cmd)
	return nil
}

// Start to start the cron scheduler
func (ss *StaticScheduler) Start() error {
	for _, job := range ss.jobs {
		job()
	}
	return nil
}
