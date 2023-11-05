package main

import (
	"fmt"
	_ "fmt"
	_ "io/fs"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
)

func listFilesAndFolders(directory string) {
	// Start by walking the specified directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		// Check if it's a directory or a file
		if !info.IsDir() {
			fmt.Println(path) // Print file path
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	// Check if a command-line argument (root directory) is provided
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	}

	// Get the root directory from the command-line argument
	directory := os.Args[1]

	// Call the listFilesAndFolders function to list files in the specified directory
	listFilesAndFolders(directory)
}
