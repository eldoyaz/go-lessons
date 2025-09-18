package less_3_2_sem

import (
	"fmt"
	"runtime"
	"sync"
	"time"
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
