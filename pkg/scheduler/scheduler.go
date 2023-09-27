package scheduler

import (
	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	AddFunc(spec string, cmd func()) error
	Start() error
}

type CronRunner struct {
	cron *cron.Cron
}

func NewCron() *CronRunner {
	return &CronRunner{cron: cron.New()}
}

func (cr *CronRunner) AddFunc(spec string, cmd func()) error {
	_, err := cr.cron.AddFunc(spec, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (cr *CronRunner) Start() error {
	cr.cron.Start()
	return nil
}

type StaticScheduler struct {
	job func()
}

func NewStaticScheduler() *StaticScheduler {
	return &StaticScheduler{}
}

func (ss *StaticScheduler) AddFunc(spec string, cmd func()) error {
	ss.job = cmd
	return nil
}

func (ss *StaticScheduler) Start() error {
	ss.job()
	return nil
}
