package scheduler

import (
	"testing"
	"time"
)

func TestGoCron_AddFunc(t *testing.T) {
	gc := NewGoCron()
	var executed bool

	gc.AddFunc("@every 1s", func() {
		executed = true
	})
	gc.Start()

	time.Sleep(2 * time.Second) // Wait 2 seconds to ensure the scheduled function runs

	if !executed {
		t.Fatalf("expected job function to be executed, but it wasn't")
	}
}

func TestGoCron_Start(t *testing.T) {
	gc := NewGoCron()
	var executed bool

	gc.AddFunc("@every 1s", func() {
		executed = true
	})

	gc.Start()

	time.Sleep(2 * time.Second) // Wait 2 seconds to ensure the scheduled function runs

	if !executed {
		t.Fatalf("expected job function to be executed, but it wasn't")
	}
}

func TestGoCron_AddFunc_Error(t *testing.T) {
	gc := NewGoCron()

	err := gc.AddFunc("INVALID_SPEC", func() {})

	if err == nil {
		t.Fatalf("expected an error for invalid cron spec, but got none")
	}
}
