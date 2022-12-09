package main

import (
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

// Create an exception list of directory names
var exclude = []string{
	"OLD",
}

func main() {
	// Get the current directory
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	// Get the list of directories in the current directory
	childDirs, err := listDirectories(currentDir, exclude)
	if err != nil {
		panic(err)
	}

	for _, d := range childDirs {
		// Run the compression in a goroutine
		wg.Add(1)
		go compress(d, currentDir, nil)
	}

	// Wait for all the goroutines to finish
	wg.Wait()
}
