package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

var (
	from, to      string
	limit, offset int64
)

const chunkS = 512 * 1024

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	wg := sync.WaitGroup{}
	progressCh := make(chan int, 5) // буфер просто чтбы не блокировать писателя (функцию копи)
	defer close(progressCh)

	completionCh := make(chan any)
	defer close(completionCh)

	progressBar(progressCh, completionCh, &wg)

	err := Copy(from, to, offset, limit, progressCh, completionCh, chunkS)
	if err != nil {
		log.Default().Println("copy failed:", err)
	}

	wg.Wait()
	fmt.Print("\n **complete** \n")
}

func progressBar(progressCh <-chan int, completionCh <-chan any, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		isFirst := true

		defer fmt.Print("]")

		for {
			select {
			case <-progressCh:
				if isFirst {
					fmt.Print("[")
					isFirst = false
				}
				fmt.Print("*")

			case <-completionCh:
				wg.Done()
				return
			}
		}
	}()
}
