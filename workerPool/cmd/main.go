package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	workerCount int
	jobCount    int
	jobChan     chan int
	wg          *sync.WaitGroup
}

func NewPool(workerCount int, jobCount int) *Pool {
	return &Pool{
		workerCount: workerCount,
		jobCount:    jobCount,
		jobChan:     make(chan int, jobCount),
		wg:          &sync.WaitGroup{},
	}
}

func (p *Pool) Start() {
	defer p.wg.Wait()

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.work(i)
	}

	defer close(p.jobChan)
	for i := 0; i < p.jobCount; i++ {
		p.jobChan <- i
	}
}

func (p *Pool) work(num int) {
	for job := range p.jobChan {
		log.Printf("worker #%d started with job #%d", num, job)

		// Имитируем работу:
		now := time.Now()
		time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)

		log.Printf("worker #%d finished with job #%d in %v", num, job, time.Since(now))
	}
	log.Printf("worker #%d stopped", num)
	p.wg.Done()
}

func main() {
	pool := NewPool(2, 5)
	pool.Start()

	//time.Sleep(2 * time.Second)
}
