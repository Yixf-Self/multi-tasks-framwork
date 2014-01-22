package main

// This is a go script illustrating how to use goroutine and channel to
// do multi-tasks.
//
// It uses the producer-consumer model.
//
//
// Author: Wei Shen <shenwei#gmail.com>
// Site:   http://shenwei.me https://github.com/shenwei356
//

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	queue      chan interface{}
	WORKER_NUM int
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s STRING [STRING...]\n", os.Args[0])
	}

	WORKER_NUM = 4

	// let the go use WORKER_NUM CPUs
	runtime.GOMAXPROCS(WORKER_NUM)

	// create a channel (queue)
	queue = make(chan interface{})

	// create workers
	for i := 1; i <= WORKER_NUM; i++ {
		fmt.Printf("Create worker %s.\n", strconv.Itoa(i))
		go worker(strconv.Itoa(i))
	}

	// channel c is a way to detect wheather the producer finished.
	c := make(chan int)

	// goroutine that dynamiclly add tasks to the queue
	go producer(c)

	//wait until all jobs being done.
	<-c
}

func producer(c chan int) {
	for _, arg := range flag.Args() {
		fmt.Printf("Enqueue: %s.\n", arg)
		queue <- arg
		time.Sleep(1 * time.Second)
	}
	fmt.Println("All tasks been sended.\nSend terimnal signal.")

	// send terimnal signal
	for i := 1; i <= WORKER_NUM; i++ {
		queue <- "STOP"
	}
	c <- 1
}

func worker(worker_name string) {
	for {
		select {
		case element := <-queue:
			if element == "STOP" {
				fmt.Printf("Worker %s stoped.\n", worker_name)
				return
			}
			result := DoSomething(element)
			fmt.Printf("Result from Worker %s: %s\n", worker_name, result)
		}
	}
}

func DoSomething(element interface{}) string {
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("I will do something with %s", element)
}
