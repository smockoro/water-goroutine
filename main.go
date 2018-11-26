package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	wordList := []string{"Apple", "Orange", "Pineapple", "Pen", "The", "A"}
	lowercaseCh := make(chan string, 10)
	uppercaseCh := make(chan string, 10)
	stopwordCh := make(chan string, 5)

}
