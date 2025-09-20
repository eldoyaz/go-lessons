package buffChannel

import (
	"context"
)

// ChanSemaphore простая реализация семафора с использованием буферизированного канала.
// Канал представляет занятые "токены" (ресурсы). Емкость канала = общему числу ресурсов.
// Изначально канал пустой. Запись в канал = занятие ресурса, чтение = освобождение.
type ChanSemaphore struct {
	tokens chan struct{} // Буферизированный канал для занятых токенов
}

// NewChanSemaphore создает новый семафор с заданным количеством ресурсов (n).
// n должно быть > 0, иначе вернет ошибку.
func NewChanSemaphore(n int64) *ChanSemaphore {
	sem := &ChanSemaphore{
		tokens: make(chan struct{}, n), // Буфер размером n для занятых ресурсов
	}
	// Канал изначально пустой
	return sem
}

// Acquire захватывает n ресурсов. Блокируется, если ресурсов недостаточно.
// Поддерживает отмену через context. Возвращает ошибку, если контекст отменен.
func (s *ChanSemaphore) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return nil // Ничего не делать для n <= 0
	}
	for i := int64(0); i < n; i++ {
		select {
		case s.tokens <- struct{}{}: // Занимаем ресурс (пишем в канал)
		case <-ctx.Done(): // Контекст отменен
			// Возвращаем уже занятые ресурсы обратно (читаем из канала)
			for j := int64(0); j < i; j++ {
				<-s.tokens
			}
			return ctx.Err()
		}
	}
	return nil
}

// TryAcquire пытается захватить n ресурсов без блокировки.
// Возвращает true, если удалось захватить все n, иначе false.
func (s *ChanSemaphore) TryAcquire(n int64) bool {
	if n <= 0 {
		return true // Ничего не делать для n <= 0
	}
	for i := int64(0); i < n; i++ {
		select {
		case s.tokens <- struct{}{}: // Занимаем ресурс (пишем в канал)
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
// чем занято (канал пустой), это ошибка (в продакшене добавьте проверки).
func (s *ChanSemaphore) Release(n int64) {
	if n <= 0 {
		return
	}
	for i := int64(0); i < n; i++ {
		<-s.tokens // Освобождаем ресурс (читаем из канала)
	}
}
