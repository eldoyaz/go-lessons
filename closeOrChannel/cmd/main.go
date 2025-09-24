package main

import (
	"fmt"
	"time"

	closeOrChannel "go-lessons/closeOrChannel/internal"
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

	<-closeOrChannel.Or(
		sig(3*time.Second),
		sig(2*time.Second),
		sig(5*time.Second),
		//sig(1*time.Hour),
	)

	fmt.Printf("done after %v\n", time.Since(start))

	time.Sleep(5 * time.Second)
	// попытка закрыть канал 2ой раз
}
