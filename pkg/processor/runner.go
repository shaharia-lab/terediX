package processor

import (
	"fmt"
)

type CronJob struct {
	spec string
	cmd  func()
}

type CronJobs struct {
	jobs []CronJob
}

type Runner struct {
	scheduler Scheduler
}

func NewRunner(scheduler Scheduler) *Runner {
	return &Runner{scheduler: scheduler}
}

func (r *Runner) AddJobs(jobs []CronJob) error {
	for _, job := range jobs {
		err := r.scheduler.AddFunc(job.spec, job.cmd)
		if err != nil {
			return fmt.Errorf("failed to add jobs to cron scheduler")
		}
	}

	return nil
}

func (r *Runner) Run() error {
	r.scheduler.Start()

	return nil
}
