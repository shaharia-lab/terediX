package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

type GoCron struct {
	cron *gocron.Scheduler
}

func NewGoCron() *GoCron {
	gc := gocron.NewScheduler(time.UTC)
	gc.TagsUnique()
	gc.WaitForScheduleAll()
	return &GoCron{cron: gocron.NewScheduler(time.UTC)}
}

func (gc *GoCron) AddFunc(spec string, cmd func()) error {
	_, err := gc.cron.CronWithSeconds(spec).SingletonMode().Do(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GoCron) Start() error {
	gc.cron.StartAsync()
	return nil
}
