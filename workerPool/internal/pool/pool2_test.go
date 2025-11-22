package pool

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestRunStopsOnErrorsLimit(t *testing.T) {
	defer goleak.VerifyNone(t)

	const (
		tasksCount = 10
		workers    = 50
		errorLimit = 2
	)

	tasks := make([]Task, tasksCount)
	for i := range tasks {
		tasks[i] = func() error {
			return DefaultTask()
		}
	}

	err := Run(tasks, workers, errorLimit)
	if !errors.Is(err, ErrErrorsLimitExceeded) {
		t.Fatalf("expected ErrErrorsLimitExceeded, got %v", err)
	}

}
