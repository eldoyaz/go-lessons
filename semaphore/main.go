package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-lessons/semaphore/buffChannel"
	"go-lessons/semaphore/freeAtomic"

	"golang.org/x/sync/semaphore"
)

// Semaphore интерфейс для семафора
type Semaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}

const (
	Semaphore_CHANNEL = "channel"
	Semaphore_ATOMIC  = "atomic"
	Semaphore_X_SYNC  = "x_sync"
)

// Тип семафора
var semType = [3]string{
	Semaphore_CHANNEL,
	Semaphore_ATOMIC,
	Semaphore_X_SYNC,
}

func main() {
	n, t := initArgs()

	var sem Semaphore
	switch *t {
	case Semaphore_CHANNEL:
		sem = buffChannel.NewChanSemaphore(*n)
	case Semaphore_ATOMIC:
		sem = freeAtomic.NewAtomic(*n)
	case Semaphore_X_SYNC:
		sem = semaphore.NewWeighted(*n)
	default:
		log.Fatalln("Не соответствует типам семафора:", semType)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Захват 3 ресурсов
	if err := sem.Acquire(ctx, 3); err != nil {
		fmt.Println("Acquire failed:", err)
		return
	}
	fmt.Println("Захвачено 3 ресурса")

	// Попытка захватить еще 3 (не должно получиться, так как осталось 2)
	if sem.TryAcquire(3) {
		fmt.Println("TryAcquire: удалось захватить 3")
	} else {
		fmt.Println("TryAcquire: не удалось (осталось 2 ресурса)")
	}

	// Освободить 2
	sem.Release(2)
	fmt.Println("Освобождено 2 ресурса")

	// Попытка захватить 4 (теперь должно получиться)
	if sem.TryAcquire(4) {
		fmt.Println("TryAcquire: удалось захватить 4")
	} else {
		fmt.Println("TryAcquire: не удалось")
	}
}

func initArgs() (n *int64, t *string) {

	n = flag.Int64("n", 0, "Макс. количество ресурсов")
	t = flag.String("t", "", "Тип семафора")
	flag.Parse() // Разбираем аргументы
	if *n <= 0 {
		log.Fatalln("Ошибка: -n должно быть положительным числом")
	}
	fmt.Printf("Получено ресурсов n = %d\n", *n)

	if *t == "" {
		log.Fatalln("Ошибка: -t должно быть строкой")
	}
	fmt.Printf("Получен тип t = %s\n", *t)

	return
}
