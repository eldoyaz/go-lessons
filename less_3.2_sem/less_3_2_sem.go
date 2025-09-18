package less_3_2_sem

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func Less32Sem() {
	maxWorkers := runtime.GOMAXPROCS(0)
	println(maxWorkers)

	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent) // буферизованный канал - семафор

	var wg sync.WaitGroup
	tasks := 10

	for i := 0; i < tasks; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			sem <- struct{}{} // acquire - блокируется, если канал полный
			fmt.Printf("Горутина %d начала работу\n", id)

			time.Sleep(1 * time.Second) // имитация работы

			fmt.Printf("Горутина %d завершила работу\n", id)
			<-sem // release - освобождаем место в семафоре
		}(i)
	}

	wg.Wait()
	fmt.Println("Все задачи завершены")
}

func Less32semV2() {
	maxConcurrent := int64(3)
	sem := semaphore.NewWeighted(maxConcurrent)

	var wg sync.WaitGroup
	tasks := 10
	ctx := context.Background()

	for i := 0; i < tasks; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			if err := sem.Acquire(ctx, 1); err != nil {
				fmt.Println("Не удалось получить семафор:", err)
				return
			}
			fmt.Printf("Горутина %d начала работу\n", id)

			time.Sleep(1 * time.Second) // имитация работы

			fmt.Printf("Горутина %d завершила работу\n", id)
			sem.Release(1)
		}(i)
	}

	wg.Wait()
	fmt.Println("Все задачи завершены")
}
