package main

import (
	"fmt"
	"go-lessons/less_3_chan"
	"time"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("took %v\n", time.Since(now))
	}(time.Now())

	less_3_chan.Less3Chan()
}
