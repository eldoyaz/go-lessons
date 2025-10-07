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

		}(ctx, ch)
	}

	return orStatus
}
