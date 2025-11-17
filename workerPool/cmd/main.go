package main

import (
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
		tasks = append(tasks, pool.DefaultTask)
	}

	err := pool.Run(tasks, 5, 5)
	if err != nil {
		log.Printf("error: %v", err)
	}

}
