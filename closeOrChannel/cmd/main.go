package main

import (
	"fmt"
	"go-lessons/internal"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	<-internal.Or(
		sig(2*time.Second),
		sig(1999*time.Millisecond),
		sig(4*time.Second),
		//sig(1*time.Hour),
	)

	fmt.Printf("done after %v\n", time.Since(start))

	time.Sleep(5 * time.Second)
	// попытка закрыть канал 2ой раз
}
