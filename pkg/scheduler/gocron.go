package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

type GoCron struct {
	cron *gocron.Scheduler
}

func NewGoCron() *GoCron {
	return &GoCron{cron: gocron.NewScheduler(time.UTC)}
}

func (gc *GoCron) AddFunc(spec string, cmd func()) error {
	_, err := gc.cron.Cron(spec).Do(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GoCron) Start() error {
	gc.cron.StartAsync()
	return nil
}
