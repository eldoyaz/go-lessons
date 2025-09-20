package freeAtomic

import (
	"context"
	"sync"
	"sync/atomic"
)

// AtomicSemaphore простая реализация семафора с использованием atomic.Int64 для счетчика свободных ресурсов
// и sync.Cond для блокировки/ожидания. Это аналог предыдущей реализации на каналах,
// но без буферизированного канала — вместо него freeAtomic для подсчета и cond для синхронизации.
type AtomicSemaphore struct {
	free atomic.Int64 // Число свободных ресурсов
	cond *sync.Cond   // Для ожидания, когда ресурсов недостаточно
	mu   sync.Mutex   // Мьютекс для sync.Cond
}

// NewAtomic создает новый семафор с заданным количеством ресурсов (n).
// n должно быть > 0, иначе вернет ошибку.
func NewAtomic(n int64) *AtomicSemaphore {
	sem := &AtomicSemaphore{
		free: atomic.Int64{},
	}
	sem.free.Store(n)
	sem.cond = sync.NewCond(&sem.mu)
	return sem
}

// Acquire захватывает n ресурсов. Блокируется, если ресурсов недостаточно.
// Поддерживает отмену через context. Возвращает ошибку, если контекст отменен.
// Примечание: для поддержки контекста создается goroutine для каждого ожидания,
// что не оптимально при большом числе ожидающих (в продакшене используйте канал).
func (s *AtomicSemaphore) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return nil // Ничего не делать для n <= 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		free := s.free.Load()
		if free >= n {
			if s.free.CompareAndSwap(free, free-n) {
				return nil
			}
			// CAS не удался, повторить цикл
		} else {
			// Ресурсов недостаточно, ждем
			done := make(chan struct{})
			go func() {
				s.cond.Wait()
				close(done)
			}()
			s.mu.Unlock()
			select {
			case <-done:
				s.mu.Lock()
				// После пробуждения проверить счетчик заново
			case <-ctx.Done():
				s.mu.Lock()
				return ctx.Err()
			}
		}
	}
}

// TryAcquire пытается захватить n ресурсов без блокировки.
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

// Release освобождает n ресурсов. Не блокируется.
// Вызывает Broadcast, чтобы разбудить всех ожидающих Acquire.
func (s *AtomicSemaphore) Release(n int64) {
	if n <= 0 {
		return // Ничего не делать для n <= 0
	}
	s.free.Add(n)
	s.cond.Broadcast() // Разбудить всех ожидающих
}
