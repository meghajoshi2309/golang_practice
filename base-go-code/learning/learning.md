# Key Learnings from the Go Duplicate File Finder

This document outlines the key programming concepts and best practices demonstrated in the Go duplicate file finder application.

---

## 1. Parsing Command-Line Arguments

In Go, you can process command-line arguments in two primary ways. The choice depends on the complexity of the arguments you need to handle.

#### Comparison: `os.Args` vs. `flag` package

| Feature           | `os.Args`                                       | `flag` Package                                                |
| ----------------- | ----------------------------------------------- | ------------------------------------------------------------- |
| **Type**          | A simple slice of strings (`[]string`).         | A powerful parser for flags (`-name=value`) and arguments.    |
| **Parsing**       | Manual. You have to parse flags and values yourself. | Automatic. It handles parsing flags, values, and types.       |
| **Usage Help**    | No built-in help text.                          | Automatically generates a help message with `-h` or `-help`.  |
| **Best For**      | Very simple cases, like reading a single filename. | CLI tools with options, flags, and sub-commands.              |

**Example from the code (using `flag`):**

The application uses the `flag` package to handle the directory path provided as an argument.

```go
// main.go

// This line parses all the registered flags from the command line.
flag.Parse() 

// flag.NArg() returns the number of arguments remaining after flags have been processed.
numberOfArgs := flag.NArg()
if numberOfArgs < 1 {
    log.Fatal("Please provide directory name...")
}

// flag.Arg(0) retrieves the first non-flag argument.
dirName := flag.Arg(0) 
```

---

## 2. Efficiently Hashing Large Files: Streaming vs. One-Shot

Calculating a file's hash is the core of this application. However, the method used can have a significant impact on performance and memory usage, especially with large files.

#### Comparison: One-Shot vs. Streaming Hashing

| Method             | One-Shot Hashing (`sha256.Sum256`)                                                              | Streaming Hashing (`io.Copy` + `sha256.New`)                                                              |
| ------------------ | ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| **How it Works**   | Reads the **entire file into memory** at once, then computes the hash.                          | Reads the file in small, fixed-size chunks and feeds each chunk into the hash function sequentially.    |
| **Memory Usage**   | **High**. A 10 GB file will consume at least 10 GB of RAM.                                      | **Low and constant**. Memory usage is minimal, regardless of the file size.                               |
| **Performance**    | Faster for very small files due to less overhead.                                               | More robust and scalable for files of any size. The standard for production-grade applications.         |
| **Risk**           | Can easily cause the application to crash with an "out of memory" error on large files.         | Safe and reliable for all file sizes.                                                                     |

**Example from the code (Streaming approach):**

The final implementation correctly uses the streaming method for robustness.

```go
// calculateHash function


func calculateHash(path string) (string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// This is one shot hashing it is good for small file
	// How it is working ?
	// 1. Pass a []byte.
	// 2. Go loads the entire file/data into memory first.
	// 3. Then it runs the SHA-256 algorithm and gives a [32]byte digest.

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil

	// But I need to think about large file also
	// Use streaming
	// h := sha256.New(); // return hash
	// if _, err := io.Copy(h, f); err != nil {
	// 	panic(err)
	// }
	// h.Sum(nil)

	// io.Copy reads the file in 4 KB chunks (default buffer).
	// Each chunk is fed into h.Write().
	// At the end, h.Sum(nil) gives final digest.
	// Memory stays small no matter if file is 10 MB or 100 GB.

	// h.Sum(nil) gives raw hash digest.
	// h.Sum(b) appends digest to b.
	// Always use h.Sum(nil) when just want hash.

}
```

---

## 3. Traversing File Systems with `filepath.Walk`

To find all files in a directory, including its subdirectories, Go provides the `filepath.Walk` function. It's an elegant and efficient way to "walk" a file tree.

**How it Works:**

`filepath.Walk` takes a root path and a "walk function." It recursively calls the walk function for every file and directory it finds.

**Example from the code:**

```go
// main.go

err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
    // This anonymous function is the "walk function".

    if err != nil {
        return err // Propagate errors up.
    }
    // Skip directories.
    if info.IsDir() {
        return nil
    }

    // Calculate hash for the file at 'path'.
    hash, err := calculateHash(path)
    // ... handle hash ...

    return nil // Returning nil tells Walk to continue.
})
```

---

## 4. Logging: `log` vs. `fmt`

The code uses both `fmt` and `log` packages. Understanding their differences is key to writing clean, maintainable Go applications.

#### Comparison: `log` vs. `fmt`

| Feature           | `fmt` Package                               | `log` Package                                                              |
| ----------------- | ------------------------------------------- | -------------------------------------------------------------------------- |
| **Default Output**| Standard Output (`stdout`).                 | Standard Error (`stderr`).                                                 |
| **Formatting**    | Just the string you provide.                | Automatically prepends each message with a timestamp (e.g., `2023/10/27 15:04:05`). |
| **Exit on Error** | No built-in functions to exit the program.  | Provides `log.Fatal()` and `log.Fatalf()` which print a message and then call `os.Exit(1)`. |
| **Use Case**      | For printing program output intended for the user (e.g., the final list of duplicates). | For printing status, warning, and error messages that are useful for debugging and monitoring. |

**Example from the code:**

```go
// Using log for a fatal error
if numberOfArgs < 1 {
    log.Fatal("Please provide directory name in args for check duplicates")
}

// Using fmt for user-facing output
fmt.Println("No duplicate files found.")
