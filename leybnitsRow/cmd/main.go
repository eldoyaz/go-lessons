package main

import (
	"context"
	"fmt"
	"time"

	"leybnitsRow/internal"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	pc := internal.NewPiCounter(31)

	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	pc.Start(ctx, cancel)
	pc.Print()

}
