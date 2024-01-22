package main

import (
	"fmt"
	"sync"
)

func main() {

	wg := &sync.WaitGroup{}
	myChan := make(chan int)

	wg.Add(2)
	//Reciver ONLY channel
	go func(wg *sync.WaitGroup, ch <-chan int) {
		fmt.Println("Recived :", <-ch)
		wg.Done()
	}(wg, myChan)

	// Send ONLY channel
	go func(wg *sync.WaitGroup, ch chan<- int) {
		ch <- 9
		wg.Done()
	}(wg, myChan)

	wg.Wait()
}
