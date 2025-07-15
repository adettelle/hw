package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currentIn := in

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

			for range in {
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
