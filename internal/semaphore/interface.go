package semaphore

import "context"

// MySemaphore интерфейс для семафора
type MySemaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}
