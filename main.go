package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// cleanDir deletes all files and subfolders in the given directory and logs the results.
func cleanDir(dirPath string) error {
	// Read the directory contents
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("reading directory %s failed: %v", dirPath, err)
	}

	// Iterate through each item and attempt to remove it.
	for _, item := range items {
		fullPath := filepath.Join(dirPath, item.Name())
		err := os.RemoveAll(fullPath)
		if err != nil {
			fmt.Printf("Failed to remove %s: %v\n", fullPath, err)
		} else {
			fmt.Printf("Successfully removed %s\n", fullPath)
		}
	}
	return nil
}

func main() {
	// Retrieve the current user's TEMP directory dynamically.
	userTemp := os.Getenv("TEMP")
	if userTemp == "" {
		// Fallback: derive from the user's home directory.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error retrieving user home directory: %v\n", err)
			return
		}
		userTemp = filepath.Join(homeDir, "AppData", "Local", "Temp")
	}

	// List of directories to clean.
	tempDirs := []string{
		userTemp,            // Current user's TEMP directory
		`C:\Windows\Temp`,   // Windows system temporary directory
		`C:\Windows\Prefetch`, // Windows prefetch directory
	}

	fmt.Println("Starting cleanup process...")

	// Iterate over each directory and perform the cleanup.
	for _, dir := range tempDirs {
		fmt.Printf("\nCleaning directory: %s\n", dir)
		if err := cleanDir(dir); err != nil {
			fmt.Printf("Error cleaning %s: %v\n", dir, err)
		}
	}

	fmt.Println("\nCleanup complete!🚀")
}
