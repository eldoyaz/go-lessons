package internal

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

type PiCounter struct {
	contextTimeoutSec int
	goroutineCount    int
	goroutineSum      chan float64
	wg                *sync.WaitGroup
}

func NewPiCounter(goroutineCount, contextTimeout int) *PiCounter {
	return &PiCounter{
		contextTimeoutSec: contextTimeout,
		goroutineCount:    goroutineCount,
		goroutineSum:      make(chan float64, goroutineCount),
		wg:                &sync.WaitGroup{},
	}
}

func (p *PiCounter) calc(ctx context.Context, i int64) {
	defer p.wg.Done()

	var sum float64

	for {
		select {
		case <-ctx.Done():
			p.goroutineSum <- sum
			return
		default:
			sum += math.Pow(-1, float64(i)) / float64(2*i+1)
			i += int64(p.goroutineCount)
		}
	}

}

func (p *PiCounter) Start() {
	defer func() {
		close(p.goroutineSum)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.contextTimeoutSec)*time.Second)
	defer cancel()

	for i := int64(0); i < int64(p.goroutineCount); i++ {
		p.wg.Add(1)
		go p.calc(ctx, i)
	}
	p.wg.Wait()
}

func (p *PiCounter) Print() {
	var sum float64

	for s := range p.goroutineSum {
		log.Println("sum:", s)
		sum += s
	}

	fmt.Printf("\nsum*4 = %.9f\n", sum*4)
}
