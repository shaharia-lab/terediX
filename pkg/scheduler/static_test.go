package scheduler

import (
	"testing"
)

func TestStaticScheduler_AddFunc(t *testing.T) {
	ss := NewStaticScheduler()
	if len(ss.jobs) != 0 {
		t.Fatalf("expected initial jobs length to be 0, but got %d", len(ss.jobs))
	}

	ss.AddFunc("", func() {
		// This is a dummy function for testing
	})
	if len(ss.jobs) != 1 {
		t.Fatalf("expected jobs length to be 1 after adding a job, but got %d", len(ss.jobs))
	}
}

func TestStaticScheduler_Start(t *testing.T) {
	var executed bool

	ss := NewStaticScheduler()
	ss.AddFunc("", func() {
		executed = true
	})

	ss.Start()

	if !executed {
		t.Fatalf("expected job function to be executed, but it wasn't")
	}
}

func TestStaticScheduler_MultipleJobs(t *testing.T) {
	counter := 0

	ss := NewStaticScheduler()
	ss.AddFunc("", func() {
		counter++
	})

	ss.AddFunc("", func() {
		counter++
	})

	ss.Start()

	if counter != 2 {
		t.Fatalf("expected counter to be 2 after executing all jobs, but got %d", counter)
	}
}
