package main

import (
	"context"
	"fmt"
	"time"

	"closeOrChannel/internal"
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

	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()

	<-internal.Or(
		ctx,
		sig(2*time.Second),
		sig(1999*time.Millisecond),
		sig(10*time.Second),
	)
	cancel()

	fmt.Printf("done after %v\n", time.Since(start))

	time.Sleep(1 * time.Second)
}
