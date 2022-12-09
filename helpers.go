package main

import (
	"os"
	"path/filepath"

	hidden "github.com/tobychui/goHidden"
)

// Returns a list of directories in the given path, excluding the directories in the given array
func listDirectories(path string, exclude []string) ([]string, error) {
	// Create a slice to hold the list of directories
	directories := []string{}

	// Walk the given path and collect the directories
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the directory is hidden
		isHidden, err := hidden.IsHidden(path, false)
		if err != nil {
			return err
		}

		if isHidden {
			return nil
		}

		// Check if the directory is in the exclusion array
		if contains(exclude, info.Name()) {
			return nil
		}

		// Add the directory to the list if it is not a file
		if info.IsDir() {
			directories = append(directories, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return directories, nil
}

// Returns true if the given array contains the given string
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
