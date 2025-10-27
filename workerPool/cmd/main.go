package main

import (
	"fmt"
	"time"
	"workerPool/internal/pool"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	myPool := pool.NewPool(
		8,
		5,
		3,
	)
	myPool.Start()

	//time.Sleep(2 * time.Second)
}
