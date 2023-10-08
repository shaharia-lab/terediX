// Package scheduler helps to schedule the task
package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

// GoCron is a wrapper of gocron
type GoCron struct {
	cron *gocron.Scheduler
}

// NewGoCron create new instance of GoCron
func NewGoCron() *GoCron {
	gc := gocron.NewScheduler(time.UTC)
	gc.TagsUnique()
	gc.WaitForScheduleAll()
	return &GoCron{cron: gc}
}

// AddFunc add new function to scheduler
func (gc *GoCron) AddFunc(spec string, cmd func()) error {
	_, err := gc.cron.CronWithSeconds(spec).SingletonMode().Do(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Start the scheduler
func (gc *GoCron) Start() error {
	gc.cron.StartAsync()
	return nil
}
