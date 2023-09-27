package scheduler

import "github.com/robfig/cron/v3"

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
