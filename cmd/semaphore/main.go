package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	mySemaphore "go-lessons/internal/semaphore"

	"golang.org/x/sync/semaphore"
)

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

func (t SemaphoreType) Create(n int64) mySemaphore.MySemaphore {

	var sem mySemaphore.MySemaphore

	switch t {
	case SemaphoreTypeChannel:
		sem = mySemaphore.NewBuffChanSemaphore(n)
	case SemaphoreTypeAtomic:
		sem = mySemaphore.NewAtomicSemaphore(n)
	case SemaphoreTypeXSync:
		sem = semaphore.NewWeighted(n)
	}
	return sem
}

func initArgs() (int64, string) {

	n := flag.Int64("n", 0, "Макс. количество ресурсов")
	t := flag.String("t", "", "Тип семафора")
	flag.Parse() // Разбираем аргументы

	if *n <= 0 {
		log.Fatalln("Ошибка: -n должно быть положительным числом")
	}
	log.Printf("Получено ресурсов n = %d\n", *n)

	if *t == "" {
		log.Fatalln("Ошибка: -t должно быть строкой")
	}
	log.Printf("Получен тип t = %s\n", *t)

	semType := SemaphoreType(*t)
	if !semType.Valid() {
		log.Fatalln("Не соответствует типам семафора:", semType)
	}

	return *n, *t
}

func main() {
	n, t := initArgs()

	sem := SemaphoreType(t).Create(n)

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
