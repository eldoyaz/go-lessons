package internal

import (
	"context"
	"fmt"
	"sync"
)

func Or(ctx context.Context, channels ...<-chan interface{}) <-chan interface{} {
	orStatus := make(chan interface{})
	once := &sync.Once{}

	for _, ch := range channels {
		go func(ctx context.Context, ch <-chan interface{}) {

			for {
				select {
				case <-ch:
					fmt.Println("Канал закрыт")
					once.Do(func() {
						close(orStatus)
					})
					return
				case <-ctx.Done():
					fmt.Println("Контекст закрыт")
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
		}(ctx, ch)
	}

	return orStatus
}
