package leybnitsRow

import (
	"math"
	"sync"
)

func Less31Chan() {

}

func CalcPart(init int64, step int64, ch chan float64, wg *sync.WaitGroup) {
	var s float64

	for i := init; i >= 0; i += step {
		s += Elem(i)
		//fmt.Printf("s=%v\n", s)
	}

	ch <- s
	wg.Done()
}

func Elem(n int64) float64 {
	//fmt.Printf("-------\nn=%d\n", n)

	e := math.Pow(-1, float64(n)) / float64(2*n+1)
	//fmt.Printf("e=%v\n", e)

	return e
}
