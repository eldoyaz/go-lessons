package main

import (
	"fmt"
	"time"

	less_3_3_or "go-lessons/less_3.3_or"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	//less_3_1_chan.Less3Chan()
	//less_3_2_sem.Less32Sem_v2()
	less_3_3_or.Less33Or()
}
