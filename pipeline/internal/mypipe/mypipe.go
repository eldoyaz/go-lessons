package mypipe

type (
	In  = <-chan interface{}
	Out = In
)

type Stage func(in In) (out Out)

// ExecutePipeline запускает конкуррентный пайплайн из стейджей.
// Каждый элемент обрабатывается независимо, что обеспечивает конкуррентность.
// Пайплайн может быть остановлен через канал done.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		// Если стейджей нет, просто возвращаем входной канал
		return in
	}
	/*
	   	var out Out

	   	// Последовательно применяем каждый стейдж
	   	for _, stage := range stages {
	   		// Применяем стейдж с обработкой done и паники
	   		out = wrapStage(stage, in, done)
	   	}

	   	return out
	   }
	*/

	// Создаем входной канал с поддержкой done
	wrappedIn := wrapInputWithDone(in, done)

	// Последовательно применяем каждый стейдж
	for _, stage := range stages {
		// Применяем стейдж с обработкой done и паники
		wrappedIn = wrapStage(stage, wrappedIn, done)
	}

	return wrappedIn
}

// wrapInputWithDone оборачивает входной канал для обработки done
func wrapInputWithDone(in In, done In) Out {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}
	}()

	return out
}

// wrapStage оборачивает стейдж для обработки done канала и паники
func wrapStage(stage Stage, in In, done In) Out {
	out := make(chan interface{})

	go func() {
		defer close(out)

		// Вызываем стейдж с защитой от паники
		var stageOut Out
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Паника в стейдже - закрываем выходной канал
					stageOut = nil
				}
			}()
			stageOut = stage(in)
		}()

		// Если стейдж упал с паникой, просто закрываем выходной канал
		if stageOut == nil {
			return
		}

		// Пересылаем данные из стейджа в выходной канал с обработкой done
		for {
			select {
			case <-done:
				return
			case val, ok := <-stageOut:
				if !ok {
					// Выходной канал стейджа закрыт
					return
				}
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}
	}()

	return out
}
