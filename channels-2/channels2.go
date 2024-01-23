package main

import (
	"fmt"
	"time"
)

type T struct {
	i byte
	b bool
}

func send(i int, ch chan<- *T) {
	t := &T{i: byte(i)}
	ch <- t
	t.b = true
}

func main() {
	vs := make([]T, 5)
	ch := make(chan *T, 5)

	for i := range vs {
		go send(i, ch)
	}

	time.Sleep(1 * time.Second)

	for i := range vs {
		vs[i] = *<-ch
	}
	for _, v := range vs {
		fmt.Println(v)
	}
}
