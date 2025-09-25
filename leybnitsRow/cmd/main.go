package main

import (
	"fmt"
	"sync"

	"go-lessons/leybnitsRow/internal"
)

const Limit = 10000000 // 10 млн.

func main() {

	var i int64
	var goroutineCount int64 = 100
	sumCh := make(chan float64, goroutineCount)

	wg := &sync.WaitGroup{}
	for ; i < goroutineCount; i++ {
		wg.Add(1)
		go leybnitsRow.CalcPart(i, goroutineCount, sumCh, wg)
	}
	wg.Wait()
	close(sumCh)

	var sum float64

	for s := range sumCh {
		sum += s
	}

	fmt.Printf("\ns*4=%.9f\n", sum*4)
}
