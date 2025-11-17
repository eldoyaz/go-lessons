package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"workerPool/internal/pool"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	taskCount := 10
	tasks := make([]pool.Task, 0, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks = append(tasks, task)
	}

	err := pool.Run(tasks, 5, 5)
	if err != nil {
		log.Printf("error: %v", err)
	}

}

var ErrTaskStopped = errors.New("task stopped")

func task() error {

	log.Printf("task started")

	// Имитируем ошибку
	rnd := rand.Int63n(5000)
	if rnd < 3000 {
		log.Print("task finished with ERR")
		return ErrTaskStopped
	}

	// Имитируем работу:
	now := time.Now()
	time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)
	log.Printf("task finished in %v", time.Since(now))

	return nil
}
