package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Есть функция, работающая неопределённо долго и возвращающая число.
// Её тело нельзя изменять (представим, что внутри сетевой запрос).
func unpredictableFunc() int64 {
	rnd := rand.Int63n(5000)
	time.Sleep(time.Duration(rnd) * time.Millisecond)

	return rnd
}

// Нужно изменить функцию обёртку, которая будет работать с заданным таймаутом (например, 1 секунду).
// Если "длинная" функция отработала за это время - отлично, возвращаем результат.
// Если нет - возвращаем ошибку. Результат работы в этом случае нам не важен.
//
// Дополнительно нужно измерить, сколько выполнялась эта функция (просто вывести в лог).
// Сигнатуру функции обёртки менять можно.
func predictableFunc(ctx context.Context) (int64, error) {

	resultChan := make(chan int64)
	//defer close(resultChan)

	go func() {
		defer close(resultChan)
		defer fmt.Println("goroutine done") // @todo почему ИНОГДА эта строчка выводится РАНЬШЕ ?

		select {
		case resultChan <- unpredictableFunc():
			println("unpredictableFunc done")
		case <-ctx.Done():
			println("ctx done")
		}
	}()

	select {
	case res := <-resultChan:
		return res, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	fmt.Println("started")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println(predictableFunc(ctx))
	fmt.Println(predictableFunc(ctx))
	fmt.Println(predictableFunc(ctx))

	time.Sleep(3 * time.Second)
}
