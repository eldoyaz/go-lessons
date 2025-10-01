package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type PiCounter struct {
	goroutineCount int
	goroutineSum   chan float64
	wg             *sync.WaitGroup
}

func NewPiCounter(goroutineCount int) *PiCounter {
	return &PiCounter{
		goroutineCount: goroutineCount,
		goroutineSum:   make(chan float64, goroutineCount),
		wg:             &sync.WaitGroup{},
	}
}

func (p *PiCounter) calc(ctx context.Context, i int) {
	defer p.wg.Done()

	var sum float64

	for {
		select {
		case <-ctx.Done():
			p.goroutineSum <- sum
			return
		default:
			sum += math.Pow(-1, float64(i)) / float64(2*i+1)
			i += p.goroutineCount
		}
	}

}

func (p *PiCounter) start() {
	defer func() {
		close(p.goroutineSum)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < p.goroutineCount; i++ {
		p.wg.Add(1)
		go p.calc(ctx, i)
	}
	p.wg.Wait()
}

func (p *PiCounter) print() {
	var sum float64

	for s := range p.goroutineSum {
		sum += s
	}

	fmt.Printf("\ns*4=%.9f\n", sum*4)
}

func main() {

	pc := NewPiCounter(200)
	pc.start()
	pc.print()

	/*
		var goroutineCount int64 = 100
		sumCh := make(chan float64, goroutineCount)

		wg := &sync.WaitGroup{}
		for i := int64(0); i < goroutineCount; i++ {
			wg.Add(1)
			go internal.CalcPart(i, goroutineCount, sumCh, wg)
		}
		wg.Wait()
		close(sumCh)

		var sum float64

		for s := range sumCh {
			sum += s
		}

		fmt.Printf("\ns*4=%.9f\n", sum*4)
	*/

}
