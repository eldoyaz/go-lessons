package main

import (
	"fmt"
	"time"

	"leybnitsRow/internal"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	pc := internal.NewPiCounter(3, 9)
	pc.Start()
	pc.Print()

}
