package main

import (
	"fmt"
	"time"

	less_3_2_sem "go-lessons/less_3.2_sem"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	//less_3_1_chan.Less31Chan()
	less_3_2_sem.Less32Sem()
	//less_3_3_or.Less33Or()
}
