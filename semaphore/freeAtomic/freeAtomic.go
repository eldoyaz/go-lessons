package freeAtomic

import (
	"context"
	"sync/atomic"
)

// AtomicSemaphore простая реализация семафора с использованием atomic.Int64 для счетчика свободных ресурсов
// и sync.Cond для блокировки/ожидания. Это аналог предыдущей реализации на каналах,
// но без буферизированного канала — вместо него freeAtomic для подсчета и cond для синхронизации.
type AtomicSemaphore struct {
	free atomic.Int64 // Число свободных ресурсов
}

// NewAtomic создает новый семафор с заданным количеством ресурсов (n).
func NewAtomic(n int64) *AtomicSemaphore {
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
			if m := s.free.Add(-n); m >= 0 {
				return nil
			}
		} else {
			// Ресурсов недостаточно, ждем
			select {
			case <-ctx.Done():
				return ctx.Err()
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
	for {
		free := s.free.Load()
		if free >= n {
			if s.free.CompareAndSwap(free, free-n) {
				return true
			}
			// CAS не удался, повторить
		} else {
			return false
		}
	}
}

// Release освобождает n ресурсов.
func (s *AtomicSemaphore) Release(n int64) {
	if n <= 0 {
		return // Ничего не делать для n <= 0
	}
	s.free.Add(n)
}
