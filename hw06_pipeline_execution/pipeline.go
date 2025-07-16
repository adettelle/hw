package hw06pipelineexecution

import "log"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currentIn := createClosableChan(in, done)

	for _, stage := range stages {
		wrapped := createClosableChan(currentIn, done)
		currentIn = stage(wrapped)
	}

	return currentIn
}

func createClosableChan(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)

			for val := range in {
				log.Println(val)
			}
		}()

		for {
			select {
			case _, ok := <-done:
				if !ok {
					return
				}
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()

	return out
}
