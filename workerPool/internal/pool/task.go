package pool

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

var ErrTaskStopped = errors.New("task stopped")

func DefaultTask() (err error) {
	defer func(now time.Time) {
		if err == nil {
			time.Sleep(time.Duration(rand.Int63n(5000)) * time.Millisecond)
			log.Printf("task finished in %v", time.Since(now))
		}
	}(time.Now())

	log.Printf("task started")

	rnd := rand.Int63n(5000)
	if rnd < 1000 {
		log.Print("task finished with ERR")
		err = ErrTaskStopped
	}

	return err
}
