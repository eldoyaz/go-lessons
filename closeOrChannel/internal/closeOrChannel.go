package internal

import (
	"context"
	"fmt"
)

func Or(channels ...<-chan interface{}) <-chan interface{} {
	orStatus := make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	for _, ch := range channels {
		go func(ch <-chan interface{}) {

			for {
				select {
				case <-ch:
					fmt.Println("Канал закрыт")
					close(orStatus)
					cancel()
					//return
				case <-ctx.Done():
					return
				}
			}

			//Проверяем статус
			//_, ok := <-ch
			//if !ok {
			//	fmt.Println("Канал закрыт")
			//
			//	_, ok1 := <-orStatus
			//	if !ok1 {
			//		fmt.Println("Канал1 закрыт")
			//		return
			//	}
			//	close(orStatus)
			//}
		}(ch)
	}

	return orStatus
}
