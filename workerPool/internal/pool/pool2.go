package pool

import (
	"errors"
	"log"
	"sync"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m < 0 {
		return nil
	}

	taskChan := make(chan Task, len(tasks))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errorCount := 0
	shouldStop := false

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// Проверяем, нужно ли остановиться
				mu.Lock()
				if shouldStop {
					log.Print("task not started: worker stopped")
					mu.Unlock()
					return
				}
				mu.Unlock()

				// Выполняем задачу
				if err := task(); err != nil {
					mu.Lock()
					errorCount++
					// Проверяем, не превышен ли лимит ошибок
					if errorCount >= m {
						shouldStop = true
						mu.Unlock()
						return
					}
					mu.Unlock()
				}
			}
		}()
	}

	// Отправляем задачи в канал
	for _, task := range tasks {
		mu.Lock()
		if shouldStop {
			mu.Unlock()
			break
		}
		mu.Unlock()
		taskChan <- task
		time.Sleep(500 * time.Millisecond)
	}

	// Закрываем канал задач
	close(taskChan)

	// Ждем завершения всех воркеров
	wg.Wait()

	// Проверяем, был ли превышен лимит ошибок
	//mu.Lock()
	//defer mu.Unlock()
	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
