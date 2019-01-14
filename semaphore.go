package main

import (
	"fmt"
	"time"

	"github.com/cipepser/go-example/semaphore"
)

func main() {
	tickets, timeout := 10, 6*time.Second
	s := semaphore.NewInterface(tickets, timeout)

	for i := 0; i <= 100; i++ {
		if err := s.Acquire(); err != nil {
			panic(err)
		}

		go func(i int) {
			doHeavyProcess(i)

			if err := s.Release(); err != nil {
				panic(err)
			}
		}(i)
	}
}

func doHeavyProcess(i int) {
	fmt.Printf("process[%d] starts\n", i)
	time.Sleep(7 * time.Second)
	fmt.Printf("process[%d] ends\n", i)
}
