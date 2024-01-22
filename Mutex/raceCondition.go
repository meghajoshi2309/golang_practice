package main

import (
	"fmt"
	"sync"
)

func main() {

	wg := &sync.WaitGroup{}
	mut := &sync.Mutex{}

	var score = []int{0}

	wg.Add(3)
	go func(wg *sync.WaitGroup, m *sync.Mutex) {
		fmt.Println("One Race")
		m.Lock()
		score = append(score, 1)
		m.Unlock()
		wg.Done()
	}(wg, mut)
	go func(wg *sync.WaitGroup, m *sync.Mutex) {
		fmt.Println("Two Race")
		m.Lock()
		score = append(score, 2)
		m.Unlock()

		wg.Done()
	}(wg, mut)
	go func(wg *sync.WaitGroup, m *sync.Mutex) {
		fmt.Println("Three Race")
		m.Lock()
		score = append(score, 3)
		m.Unlock()
		wg.Done()
	}(wg, mut)

	wg.Wait()

	fmt.Println(score)

}

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {

// 	wg := &sync.WaitGroup{}
// 	// mut := &sync.Mutex{}

// 	var score = []int{0}

// 	wg.Add(3)
// 	go func(wg *sync.WaitGroup) {
// 		fmt.Println("One Race")
// 		// m.Lock()
// 		score = append(score, 1)
// 		// m.Unlock()
// 		wg.Done()
// 	}(wg)
// 	go func(wg *sync.WaitGroup) {
// 		fmt.Println("Two Race")
// 		// m.Lock()
// 		score = append(score, 2)
// 		// m.Unlock()

// 		wg.Done()
// 	}(wg)
// 	go func(wg *sync.WaitGroup) {
// 		fmt.Println("Three Race")
// 		// m.Lock()
// 		score = append(score, 3)
// 		// m.Unlock()
// 		wg.Done()
// 	}(wg)

// 	wg.Wait()

// 	fmt.Println(score)

// }
