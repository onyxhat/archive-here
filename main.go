package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	hidden "github.com/tobychui/goHidden"
)

var wg sync.WaitGroup

func main() {
	// Get the current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	// Get the list of directories in the current directory
	dirs, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		panic(err)
	}

	// Create an exception list of directory names
	exceptions := []string{
		"OLD",
	}

	for _, d := range dirs {
		// Check if the directory is in the exception list
		found := false
		for _, e := range exceptions {
			if strings.HasSuffix(d, e) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		// Check if the directory is hidden
		isHidden, err := hidden.IsHidden(d, false)
		if err != nil {
			panic(err)
		}

		if isHidden {
			continue
		}

		// Skip files in the local directory
		info, err := os.Stat(d)
		if err != nil {
			panic(err)
		}

		if !info.IsDir() {
			continue
		}

		// Run the compression in a goroutine
		wg.Add(1)
		go compress(d)
	}

	// Wait for all the goroutines to finish
	wg.Wait()
}

func compress(d string) {
	defer wg.Done()

	// Create the zip file
	zipfile, err := os.Create(d + ".zip")
	if err != nil {
		panic(err)
	}
	defer zipfile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	// Compress the files in the directory
	err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		// Create a new zip file for the current file
		zipfile, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		// Open the file to compress
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file to the zip file
		_, err = io.Copy(zipfile, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(d)
	if err != nil {
		panic(err)
	}
}
