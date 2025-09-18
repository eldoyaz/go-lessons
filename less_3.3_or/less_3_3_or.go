package less_3_3_or

import (
	"fmt"
	"time"
)

//var or func(channels ...<-chan interface{}) <-chan interface{}

func Less33Or() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(1*time.Minute),
		sig(3*time.Second),
		sig(5*time.Second),
		//sig(1*time.Hour),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	orStatus := make(chan interface{})

	for _, ch := range channels {
		go func() {
			// Проверяем статус
			value, ok := <-ch
			if !ok {
				fmt.Println("Канал закрыт")
				//orStatus <- struct{}{}
				close(orStatus)
				return
			} else {
				fmt.Printf("Получено значение: %v\n", value)
			}
		}()
	}

	//loop:
	//	for {
	//		select {
	//		case <-orStatus:
	//			println("orStatus closed")
	//			time.Sleep(100 * time.Millisecond)
	//			close(orStatus)
	//			break loop
	//		}
	//	}

	return orStatus
}
