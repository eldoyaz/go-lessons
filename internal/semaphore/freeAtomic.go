package semaphore

import (
	"context"
	"sync/atomic"
)

// AtomicSemaphore простая реализация семафора с использованием atomic.Int64 для счетчика свободных ресурсов
type AtomicSemaphore struct {
	free atomic.Int64 // Число свободных ресурсов
}

// NewAtomicSemaphore создает новый семафор с заданным количеством ресурсов (n).
func NewAtomicSemaphore(n int64) *AtomicSemaphore {
	sem := &AtomicSemaphore{
		free: atomic.Int64{},
	}
	sem.free.Store(n)

	return sem
}

// Acquire захватывает n ресурсов.
func (s *AtomicSemaphore) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return nil // Ничего не делать для n <= 0
	}

	for {
		free := s.free.Load()
		if free >= n {
			// @todo: брать по 1 ресурсу
			if s.free.CompareAndSwap(free, free-n) {
				return nil
			}
		} else {
			// Ресурсов недостаточно, ждем
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				continue
			}
		}
	}
}

// TryAcquire пытается захватить n ресурсов.
// Возвращает true, если удалось захватить все n, иначе false.
func (s *AtomicSemaphore) TryAcquire(n int64) bool {
	if n <= 0 {
		return true // Ничего не делать для n <= 0
	}

	free := s.free.Load()
	if free >= n {
		return s.free.CompareAndSwap(free, free-n)
	}

	return false
}

// Release освобождает n ресурсов.
func (s *AtomicSemaphore) Release(n int64) {
	if n <= 0 {
		return // Ничего не делать для n <= 0
	}

	s.free.Add(n)
}
