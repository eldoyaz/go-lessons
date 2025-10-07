package internal

import "golang.org/x/sync/semaphore"

type SemaphoreType string

const (
	SemaphoreTypeChannel SemaphoreType = "channel"
	SemaphoreTypeAtomic  SemaphoreType = "atomic"
	SemaphoreTypeXSync   SemaphoreType = "x_sync"
)

func (t SemaphoreType) Valid() bool {
	return t == SemaphoreTypeChannel ||
		t == SemaphoreTypeAtomic ||
		t == SemaphoreTypeXSync
}

func (t SemaphoreType) Create(n int64) MySemaphore {

	var sem MySemaphore

	switch t {
	case SemaphoreTypeChannel:
		sem = NewBuffChanSemaphore(n)
	case SemaphoreTypeAtomic:
		sem = NewAtomicSemaphore(n)
	case SemaphoreTypeXSync:
		sem = semaphore.NewWeighted(n)
	}
	return sem
}
