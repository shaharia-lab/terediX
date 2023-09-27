package scheduler

import (
	"github.com/robfig/cron/v3"
)

// Scheduler interface
type Scheduler interface {
	AddFunc(spec string, cmd func()) error
	Start() error
}

// CronRunner definition
type CronRunner struct {
	cron *cron.Cron
}

// NewCron to create new cron runner
func NewCron() *CronRunner {
	return &CronRunner{cron: cron.New()}
}

// AddFunc to add new function to cron scheduler
func (cr *CronRunner) AddFunc(spec string, cmd func()) error {
	_, err := cr.cron.AddFunc(spec, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Start to start the cron scheduler
func (cr *CronRunner) Start() error {
	cr.cron.Start()
	return nil
}

// StaticScheduler definition
type StaticScheduler struct {
	job func()
}

// NewStaticScheduler to create new static scheduler
func NewStaticScheduler() *StaticScheduler {
	return &StaticScheduler{}
}

// AddFunc to add new function to cron scheduler
func (ss *StaticScheduler) AddFunc(spec string, cmd func()) error {
	ss.job = cmd
	return nil
}

// Start to start the cron scheduler
func (ss *StaticScheduler) Start() error {
	ss.job()
	return nil
}
