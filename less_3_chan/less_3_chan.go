package less_3_chan

import (
	"fmt"
	"math"
	"sync"
)

var sum float64

const limit = 1000000 // 1млн.

func Less3Chan() {
	goroutineCount := 100
	sumCh := make(chan float64, goroutineCount)

	wg := &sync.WaitGroup{}

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go calcPart(i, goroutineCount, sumCh, wg)
	}
	wg.Wait()
	close(sumCh)

	for s := range sumCh {
		sum += s
	}

	fmt.Printf("\ns*4=%.9f\n", sum*4)
}

func calcPart(init, step int, ch chan float64, wg *sync.WaitGroup) {
	var s float64

	for i := init; i < limit; i += step {
		s += elem(i)
		//fmt.Printf("s=%v\n", s)
	}

	ch <- s
	wg.Done()
}

func elem(n int) float64 {
	//fmt.Printf("-------\nn=%d\n", n)

	e := math.Pow(-1, float64(n)) / float64(2*n+1)
	//fmt.Printf("e=%v\n", e)

	return e
}
