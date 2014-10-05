package main

// This is a go script illustrating how to use goroutine and channel to
// do multi-tasks.
//
// Author: Wei Shen <shenwei356#gmail.com>
// Site:   http://shenwei.me https://github.com/shenwei356
//

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	queue      chan string
	WORKER_NUM int
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s STRING [STRING...]\n", os.Args[0])
		os.Exit(1)
	}

	WORKER_NUM = 4

	// let the go use WORKER_NUM CPUs
	runtime.GOMAXPROCS(WORKER_NUM)

	// create a channel (queue)
	queue = make(chan string)

	// Producer
	go func() {
		for i, arg := range os.Args {
			if i == 0 {
				continue
			}
			fmt.Printf("Enqueue: %s.\n", arg)
			queue <- arg

			time.Sleep(1 * time.Second)
		}

		fmt.Println("All tasks been sended.\n")
		close(queue)
	}()

	// Worker
	var wg sync.WaitGroup
	tokens := make(chan int, WORKER_NUM)

	for arg := range queue {
		tokens <- 1
		wg.Add(1)

		go func(arg string) {
			defer func() {
				wg.Done()
				<-tokens
			}()

			result := DoSomething(arg)
			fmt.Printf("Result: %s\n", result)
		}(arg)
	}
	wg.Wait()
}

func DoSomething(arg string) string {
	fmt.Printf("Start to do something with %s\n", arg)
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("result of %s", arg)
}
