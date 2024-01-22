package main

import (
	"fmt"
	"net/http"
	"sync"
)

var signals []string

var wg sync.WaitGroup
var mute sync.Mutex

func main() {
	websiteList := []string{
		"https://google.com",
		"https://go.dev",
		"https://github.com",
		"https://fb.com",
		"https://lco.dev",
		"https://gorm.io/docs/index.html",
		"https://www.youtube.com/",
	}
	for _, website := range websiteList {
		wg.Add(1)
		go getStatuscode(website)
	}

	wg.Wait()
	fmt.Println("Signals : ", signals)

}

func getStatuscode(endpoint string) {

	defer wg.Done()

	mute.Lock()
	signals = append(signals, endpoint)
	mute.Unlock()

	res, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d status of %s\n", res.StatusCode, endpoint)
	}
}
