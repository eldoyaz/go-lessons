package mySemaphore

import (
	"context"
)

// BuffChanSemaphore простая реализация семафора с использованием буферизированного канала.
// Канал представляет занятые "токены" (ресурсы). Емкость канала = общему числу ресурсов.
// Изначально канал пустой. Запись в канал = занятие ресурса, чтение = освобождение.
type BuffChanSemaphore struct {
	tokens chan struct{} // Буферизированный канал для занятых токенов
}

// NewBuffChanSemaphore NewSemaphore создает новый семафор с заданным количеством ресурсов (n).
func NewBuffChanSemaphore(n int64) *BuffChanSemaphore {
	sem := &BuffChanSemaphore{
		tokens: make(chan struct{}, n),
	}
	return sem
}

// Acquire захватывает n ресурсов.
func (s *BuffChanSemaphore) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return nil // Ничего не делать для n <= 0
	}
	for i := int64(0); i < n; i++ {
		select {
		case s.tokens <- struct{}{}: // Блокируется, если ресурсов недостаточно
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// TryAcquire пытается захватить n ресурсов без блокировки.
// Возвращает true, если удалось захватить все n, иначе false.
func (s *BuffChanSemaphore) TryAcquire(n int64) bool {
	if n <= 0 {
		return true // Ничего не делать для n <= 0
	}

	//len(s.tokens) // длина канала @todo
	// если меньше n , то начинаем писать

	for i := int64(0); i < n; i++ {
		select {
		case s.tokens <- struct{}{}:
		default: // Нет доступных ресурсов (канал полон)
			// Возвращаем уже занятые ресурсы обратно (читаем из канала)
			for j := int64(0); j < i; j++ {
				<-s.tokens
			}
			return false
		}
	}
	return true
}

// Release освобождает n ресурсов. Не блокируется, но если освобождается больше,
// чем занято (канал пустой), это ошибка.
func (s *BuffChanSemaphore) Release(n int64) {
	if n <= 0 {
		return
	}
	for i := int64(0); i < n; i++ {
		select {
		case <-s.tokens: // Освобождаем ресурс (читаем из канала)
		default:
			return
		}
	}
}
