package internal

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type PiCounter struct {
	contextTimeoutSec int
	goroutineCount    int
	goroutineSum      chan float64
	signChan          chan os.Signal
	wg                *sync.WaitGroup
}

func NewPiCounter(goroutineCount, contextTimeout int) *PiCounter {
	return &PiCounter{
		contextTimeoutSec: contextTimeout,
		goroutineCount:    goroutineCount,
		goroutineSum:      make(chan float64, goroutineCount),
		signChan:          make(chan os.Signal, 1),
		wg:                &sync.WaitGroup{},
	}
}

func (p *PiCounter) calc(ctx context.Context, cancel context.CancelFunc, i int64) {
	defer p.wg.Done()

	var sum float64

	for {
		select {
		case s := <-p.signChan:
			cancel()
			p.goroutineSum <- sum
			log.Println("signal received:", s.String())
			return
		case <-ctx.Done():
			p.goroutineSum <- sum
			log.Println("context done")
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

	signal.Notify(p.signChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.contextTimeoutSec)*time.Second)
	defer cancel()

	for i := int64(0); i < int64(p.goroutineCount); i++ {
		p.wg.Add(1)
		go p.calc(ctx, cancel, i)
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
