package main

import (
	"fmt"
	"time"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	//less_3_1_chan.Less31Chan()
	//less_3_2_sem.Less32Sem()
	//less_3_3_or.Less33Or()

}
