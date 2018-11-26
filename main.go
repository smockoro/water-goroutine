package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	select_forloop()
}

func wait_group() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("No.1 goroutine sleep")
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("No.2 goroutine sleep")
		time.Sleep(1 * time.Second)
	}()

	wg.Wait()
	fmt.Println("All goroutine Done")
}

func wait_group_counter() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello, %d\n", id)
	}
	const num_Greeter = 5
	wg.Add(num_Greeter)
	for i := 0; i < num_Greeter; i++ {
		go hello(&wg, i+1)
	}

	wg.Wait()
}

func channel_hello() {
	stringCh := make(chan string)
	go func() {
		stringCh <- "Hello World."
	}()
	fmt.Println(<-stringCh)

	var writeCh chan<- string // 書き込み専用チャネルの宣言
	var readCh <-chan string  // 読み込み専用チャネルの宣言
	writeCh = stringCh
	readCh = stringCh

	go func() {
		writeCh <- "Write Channel used."
	}()

	fmt.Println(<-readCh)
}

func channel_close_for() {
	intCh := make(chan int)

	go func() {
		defer close(intCh)
		for i := 0; i < 5; i++ {
			intCh <- i
		}
	}()

	// 第二戻り値にチャネルのクローズが入っているけどfor文がうまいこと処理してくれる
	for i := range intCh {
		fmt.Println(i)
	}
}

func channel_close_multi_goroutine() {
	beganCh := make(chan interface{})
	for i := 0; i < 5; i++ {
		wg.Add(11)
		go func(i int) {
			defer wg.Done()
			//<-beganCh
			fmt.Println(<-beganCh)
			fmt.Printf("Goroutine No.%d started\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutine...")
	close(beganCh) // ここでクローズするとゴルーチンが止まる？
	wg.Wait()
}

func channel_buffer() {
	intCh := make(chan int, 4)
	i := 0

	// 1秒で1文字入れるgoroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(1 * time.Second)
			i++
			intCh <- i
		}
	}()

	// 2秒で1文字出力するgoroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(2 * time.Second)
			fmt.Println(<-intCh)
		}
	}()

	wg.Wait()
}

func select_multi_chan() {
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)

	var c1counter, c2coutner int

	for i := 0; i < 1000; i++ {
		select {
		case <-ch1:
			c1counter++
		case <-ch2:
			c2coutner++
		}
	}

	// カウンタの回数はほぼ均等になる。条件が同じチャネルを2つとも処理すると、
	// どうやら均等に処理される。
	fmt.Printf("c1counter: %d, c2coutner: %d \n", c1counter, c2coutner)
}

func select_forloop() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done) // 5秒後にチャネルクローズでdoneに情報が渡る
	}()

	loopCounter := 0
loop: // 中のbreakは無名だとselectからしか出れない、forから抜けるためにはラベル付しておく
	for {
		select {
		case <-done:
			break loop // 5秒立つとここからぬける。
		default:
		}

		loopCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("loop count: %d\n", loopCounter)
}
