package closeOrChannel

import (
	"context"
	"fmt"
)

func Or(channels ...<-chan interface{}) <-chan interface{} {
	orStatus := make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	for _, ch := range channels {
		go func(ctx context.Context, ch <-chan interface{}) {

			for {
				select {
				case orStatus <- <-ch:
					fmt.Println("Канал закрыт")
					cancel()
				case <-ctx.Done():
					return
				}
			}

			// Проверяем статус
			//_, ok := <-ch
			//if !ok {
			//	fmt.Println("Канал закрыт")
			//	close(orStatus)
			//}
		}(ctx, ch)
	}

	return orStatus
}
