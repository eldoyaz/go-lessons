package main

import (
	"fmt"
	"mypipe/internal/mypipe"
	"time"
)

func main() {
	defer func(now time.Time) {
		fmt.Printf("main() took %v\n", time.Since(now))
	}(time.Now())

	fmt.Println("=== Пример 1: Базовый пайплайн ===")
	example1()

	fmt.Println("\n=== Пример 2: Проверка конкуррентности ===")
	example2()

	fmt.Println("\n=== Пример 3: Остановка через done канал ===")
	example3()
}

// example1 демонстрирует базовую работу пайплайна
func example1() {
	// Создаем входной канал
	in := make(chan interface{})
	done := make(chan interface{})

	// Определяем стейджи
	stage1 := func(in mypipe.In) mypipe.Out {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				time.Sleep(2 * time.Second)
				// Имитация работы - умножение на 2
				result := val.(int) * 2
				out <- result
			}
		}()
		return out
	}

	stage2 := func(in mypipe.In) mypipe.Out {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				time.Sleep(2 * time.Second)
				// Имитация работы - добавление 10
				result := val.(int) + 10
				out <- result
			}
		}()
		return out
	}

	stage3 := func(in mypipe.In) mypipe.Out {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				time.Sleep(2 * time.Second)
				// Имитация работы - умножение на 3
				result := val.(int) * 3
				out <- result
			}
		}()
		return out
	}

	// Запускаем пайплайн
	out := mypipe.ExecutePipeline(in, done, stage1, stage2, stage3)

	// Отправляем данные
	go func() {
		defer close(in)
		for i := 1; i <= 5; i++ {
			fmt.Printf("Отправляем: %d\n", i)
			in <- i
		}
	}()

	// Читаем результаты
	for result := range out {
		fmt.Printf("Результат: %d\n", result.(int))
	}
}

// example2 демонстрирует конкуррентность пайплайна
func example2() {
	// Создаем входной канал
	in := make(chan interface{})
	done := make(chan interface{})

	// Стейдж с задержкой 100мс
	slowStage := func(duration time.Duration) mypipe.Stage {
		return func(in mypipe.In) mypipe.Out {
			out := make(chan interface{})
			go func() {
				defer close(out)
				for val := range in {
					time.Sleep(duration)
					out <- val
				}
			}()
			return out
		}
	}

	// Создаем 4 стейджа по 100мс каждый
	stages := []mypipe.Stage{
		slowStage(100 * time.Millisecond),
		slowStage(100 * time.Millisecond),
		slowStage(100 * time.Millisecond),
		slowStage(100 * time.Millisecond),
	}

	// Запускаем пайплайн
	out := mypipe.ExecutePipeline(in, done, stages...)

	// Отправляем 5 элементов
	go func() {
		defer close(in)
		for i := 1; i <= 5; i++ {
			fmt.Printf("Отправляем элемент: %d\n", i)
			in <- i
		}
	}()

	// Читаем результаты и замеряем время
	start := time.Now()
	count := 0
	for result := range out {
		count++
		fmt.Printf("Получен результат: %d (время с начала: %v)\n", result.(int), time.Since(start))
	}
	fmt.Printf("Всего обработано: %d элементов за %v\n", count, time.Since(start))
	fmt.Printf("Если бы было последовательно: 4 стейджа * 100мс * 5 элементов = 2 секунды\n")
	fmt.Printf("Конкуррентно получилось: %v (должно быть значительно быстрее)\n", time.Since(start))
}

// example3 демонстрирует остановку пайплайна через done канал
func example3() {
	// Создаем входной канал
	in := make(chan interface{})
	done := make(chan interface{})

	// Стейдж с задержкой
	slowStage := func(in mypipe.In) mypipe.Out {
		out := make(chan interface{})
		go func() {
			defer close(out)
			for val := range in {
				// Имитация работы - задержка 200мс
				time.Sleep(200 * time.Millisecond)
				out <- val
			}
		}()
		return out
	}

	// Запускаем пайплайн
	out := mypipe.ExecutePipeline(in, done, slowStage, slowStage)

	// Отправляем данные
	go func() {
		defer close(in)
		for i := 1; i <= 10; i++ {
			select {
			case <-done:
				return
			case in <- i:
				fmt.Printf("Отправлено: %d\n", i)
			}
		}
	}()

	// Останавливаем пайплайн через 500мс
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Отправляем сигнал остановки...")
		close(done)
	}()

	// Читаем результаты
	count := 0
	for result := range out {
		count++
		fmt.Printf("Получен результат: %d\n", result.(int))
	}
	fmt.Printf("Обработано элементов до остановки: %d\n", count)
}
