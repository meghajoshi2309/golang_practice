package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	startTime := time.Now()
	log.Println("Scan started at:", startTime.Format(time.RFC1123))

	// 1. Take dir name from CLI
	// 2. Check all file in dir
	// 3. Calculate hash of every file for check same content
	// 4. Compare hash
	// 5. Print duplicates

	dir := getDirectoryFromArgs()
	duplicates, err := findDuplicateFiles(dir)
	if err != nil {
		log.Fatalf("Error finding duplicates in directory %s: %v", dir, err)
	}

	printDuplicates(duplicates)

	duration := time.Since(startTime)
	log.Println("Scan finished in:", duration)
}

// 1. Take dir name from CLI
// getDirectoryFromArgs parses the command-line arguments to get the target directory.
func getDirectoryFromArgs() string {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("Please provide a directory name as an argument.")
	}
	return flag.Arg(0)
}

// findDuplicateFiles walks through the given directory, calculates hashes for each file,
// and returns a map of duplicate files, keyed by their hash.
func findDuplicateFiles(rootDir string) (map[string][]string, error) {
	hashes := make(map[string][]string)
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		hash, err := calculateFileHash(path)
		if err != nil {
			log.Printf("Warning: Could not calculate hash for %s: %v\n", path, err)
			return nil // Continue to the next file
		}
		hashes[hash] = append(hashes[hash], path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filterDuplicates(hashes), nil
}

// filterDuplicates takes a map of file hashes to file paths and returns a new map
// containing only the hashes that correspond to more than one file.
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
