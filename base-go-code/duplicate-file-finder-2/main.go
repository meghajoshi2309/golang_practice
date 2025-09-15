package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	log.Println("Scan started at:", startTime.Format(time.RFC1123))

	dir := getDirectoryFromArgs()
	duplicates, err := findDuplicateFiles(dir)
	if err != nil {
		log.Fatalf("Error finding duplicates in directory %s: %v", dir, err)
	}

	printDuplicates(duplicates)

	duration := time.Since(startTime)
	log.Println("Scan finished in:", duration)
}

func getDirectoryFromArgs() string {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("Please provide a directory name as an argument.")
	}
	return flag.Arg(0)
}

func findDuplicateFiles(rootDir string) (map[string][]string, error) {
	hashes := make(map[string][]string)

	filePathChan := make(chan string)

	waitGroup := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for path := range filePathChan {
				hash, err := calculateFileHash(path)
				if err != nil {
					log.Printf("Warning: Could not calculate hash for %s: %v\n", path, err)
					continue
				}

				// If I haven't lock it -> race condition occurs. Error is (fatal error: concurrent map writes).
				mu.Lock()
				hashes[hash] = append(hashes[hash], path)
				mu.Unlock()
			}

		}()
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filePathChan <- path
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filterDuplicates(hashes), nil
}

func filterDuplicates(hashes map[string][]string) map[string][]string {
	duplicates := make(map[string][]string)
	for hash, files := range hashes {
		if len(files) > 1 {
			duplicates[hash] = files
		}
	}
	return duplicates
}

func printDuplicates(duplicates map[string][]string) {
	if len(duplicates) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	fmt.Println("Found duplicate files:")
	for hash, files := range duplicates {
		fmt.Printf("\nHash: %s\n", hash)
		for _, file := range files {
			fmt.Printf("- %s\n", file)
		}
	}
}

func calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
