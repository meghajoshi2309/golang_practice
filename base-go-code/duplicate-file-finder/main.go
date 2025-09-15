package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// main function to orchestrate the duplicate file finding and deletion process.
func main() {
	// Define a command-line flag to enable deletion of duplicates.
	delete := flag.Bool("delete", false, "Set this flag to delete duplicate files.")
	flag.Parse()

	// Check if a directory path is provided as an argument.
	if flag.NArg() < 1 {
		log.Fatal("Error: Please provide a directory path to scan.")
	}
	root := flag.Arg(0)

	// Find all duplicate files within the given directory.
	duplicates := findDuplicates(root)

	// If no duplicates are found, print a message and exit.
	if len(duplicates) == 0 {
		fmt.Println("No duplicate files found.")
		return
	}

	// Print the found duplicates.
	fmt.Println("Duplicate files found:")
	for hash, files := range duplicates {
		fmt.Printf("\nHash: %s\n", hash)
		for _, file := range files {
			fmt.Printf("- %s\n", file)
		}
	}

	// If the delete flag is set, ask for confirmation and delete the files.
	if *delete {
		fmt.Print("\nAre you sure you want to delete the duplicate files? (yes/no): ")
		var response string
		fmt.Scanln(&response) // Read user input.

		if strings.ToLower(response) == "yes" {
			deleteDuplicates(duplicates)
		} else {
			fmt.Println("Deletion cancelled.")
		}
	}
}

// findDuplicates scans a directory and returns a map of hashes to file paths for duplicate files.
func findDuplicates(root string) map[string][]string {
	hashes := make(map[string][]string)
	duplicates := make(map[string][]string)

	// Walk the directory tree.
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Calculate the hash of the file content.
		hash, err := fileHash(path)
		if err != nil {
			log.Printf("Warning: Could not calculate hash for %s: %v\n", path, err)
			return nil // Continue to the next file.
		}

		// Add the file path to the map.
		hashes[hash] = append(hashes[hash], path)
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory %s: %v\n", root, err)
	}

	// Filter for hashes that have more than one file path (i.e., duplicates).
	for hash, files := range hashes {
		if len(files) > 1 {
			duplicates[hash] = files
		}
	}

	return duplicates
}

// deleteDuplicates removes all but the first file from each set of duplicates.
func deleteDuplicates(duplicates map[string][]string) {
	for _, files := range duplicates {
		// The first file (files[0]) is kept, the rest are deleted.
		for i := 1; i < len(files); i++ {
			file := files[i]
			err := os.Remove(file)
			if err != nil {
				log.Printf("Error deleting file %s: %v\n", file, err)
			} else {
				fmt.Printf("Deleted: %s\n", file)
			}
		}
	}
	fmt.Println("\nDeletion of duplicate files complete.")
}

// fileHash calculates the SHA256 hash of a file's content.
// It uses os.ReadFile for simplicity, which reads the entire file into memory.
func fileHash(path string) (string, error) {
	// Read the entire file content.
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	// Calculate the SHA256 hash.
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}
