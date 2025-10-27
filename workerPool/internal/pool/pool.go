package pool

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	errPossibleCount int
	workerCount      int
	jobCount         int
	jobChan          chan int
	wg               *sync.WaitGroup
}

func NewPool(workerCount, jobCount, errCount int) *Pool {
	return &Pool{
		errPossibleCount: errCount,
		workerCount:      workerCount,
		jobCount:         jobCount,
		jobChan:          make(chan int, jobCount),
		wg:               &sync.WaitGroup{},
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

		// Имитируем ошибку
		rnd := rand.Int63n(5000)
		if rnd < 4000 {
			log.Printf("ERR #%d: worker #%d stopped with job #%d", p.errPossibleCount, num, job)

			if p.errPossibleCount < 0 {
				break
			} else {
				p.errPossibleCount--
				continue
			}
		}

		// Имитируем работу:
		now := time.Now()
		time.Sleep(time.Duration(rand.Int63n(15000)) * time.Millisecond)

		log.Printf("worker #%d finished with job #%d in %v", num, job, time.Since(now))
	}
	log.Printf("worker #%d stopped", num)
	p.wg.Done()
}
