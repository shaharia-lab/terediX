package scheduler

import (
	"fmt"
)

// CronJob definition
type CronJob struct {
	spec string
	cmd  func()
}

// CronJobs definition
type CronJobs struct {
	jobs []CronJob
}

// Runner definition
type Runner struct {
	scheduler Scheduler
}

// AddJobs to add jobs to cron scheduler
func (r *Runner) AddJobs(jobs []CronJob) error {
	for _, job := range jobs {
		err := r.scheduler.AddFunc(job.spec, job.cmd)
		if err != nil {
			return fmt.Errorf("failed to add jobs to cron scheduler")
		}
	}

	return nil
}

// Run to start the scheduler
func (r *Runner) Run() error {
	r.scheduler.Start()

	return nil
}
