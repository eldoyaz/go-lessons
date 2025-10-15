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

func unpredictableChan() chan int64 {
	result := make(chan int64)
	go func() {
		result <- unpredictableFunc()
	}()
	return result
}

// Нужно изменить функцию обёртку, которая будет работать с заданным таймаутом (например, 1 секунду).
// Если "длинная" функция отработала за это время - отлично, возвращаем результат.
// Если нет - возвращаем ошибку. Результат работы в этом случае нам не важен.
//
// Дополнительно нужно измерить, сколько выполнялась эта функция (просто вывести в лог).
// Сигнатуру функцию обёртки менять можно.
func predictableFunc() (int64, error) {

	result := make(chan int64)
	defer close(result)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go func(ctx context.Context) {
		select {
		case res := <-unpredictableChan():
			result <- res
		case <-ctx.Done():
			result <- 0
		}
	}(ctx)

	return <-result, nil
}

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	fmt.Println("started")
	//fmt.Println(predictableFunc())
	fmt.Println(predictableFunc_v2())
}

func predictableFunc_v2() (int64, error) {

	resultChan := make(chan int64)
	defer close(resultChan)

	go func() {
		resultChan <- unpredictableFunc()
		// @todo: надо ли прерывать эту горутину, если контекст завершился по таймауту?
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	select {
	case res := <-resultChan:
		return res, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}

}
