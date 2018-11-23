package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	go func() {
		count("sheep")
		wg.Done()
	}()
	wg.Wait()
}

func count(thing string) {
	for i := 0; i <= 5; i++ {
		fmt.Println(i, thing)
		time.Sleep(time.Millisecond * 500)
	}
}
