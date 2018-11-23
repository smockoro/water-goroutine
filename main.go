package main

import "fmt"

func main() {
	c := make(chan string)
	c <- "hello"

	msg := <-c
	fmt.Println(msg)
}
