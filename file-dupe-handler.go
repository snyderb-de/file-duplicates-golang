package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	_ "fmt"
	"io"
	_ "io/fs"
	"os"
	_ "os"
	"path/filepath"
	_ "path/filepath"
	"sort"
)

func listFilesAndFolders(directory string, fileFormat string, descending bool) map[int64][]string {
	// Create a map to store files by size
	filesBySize := make(map[int64][]string)

	// Start by walking the specified directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		// Check if it's a directory or a file
		if !info.IsDir() && (fileFormat == "" || filepath.Ext(path) == "."+fileFormat) {
			size := info.Size() // Get file size
			filesBySize[size] = append(filesBySize[size], path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	// Sort the sizes in ascending order or descending order based on user input
	sizes := make([]int64, 0, len(filesBySize))
	for size := range filesBySize {
		sizes = append(sizes, size)
	}
	if descending {
		sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })
	} else {
		sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
	}

	// Print files grouped by size
	for _, size := range sizes {
		fmt.Println(size, "bytes")
		for _, path := range filesBySize[size] {
			fmt.Println(path)
		}
	}

	return filesBySize
}

// Hash files from ListFilesAndFolders using MD5
func hashFiles(filesBySize map[int64][]string, dupeCheck bool) {
	if dupeCheck {
		for size, files := range filesBySize {
			if len(files) > 1 {
				filesByHash := make(map[string][]string)

				// Hashing and grouping files by hash
				for _, file := range files {
					f, err := os.Open(file)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}
					defer f.Close()

					h := md5.New()
					if _, err := io.Copy(h, f); err != nil {
						fmt.Println("Error:", err)
						continue
					}

					hash := hex.EncodeToString(h.Sum(nil))
					filesByHash[hash] = append(filesByHash[hash], file)
				}

				// Check and print only if there are duplicates
				for hash, groupedFiles := range filesByHash {
					if len(groupedFiles) > 1 {
						fmt.Printf("\n%d bytes\n", size)
						fmt.Printf("Hash: %s\n", hash)
						for i, file := range groupedFiles {
							fmt.Printf("%d. %s\n", i+1, file)
						}
					}
				}
			}
		}
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

	// Read user input for file format
	fmt.Print("Enter file format: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fileFormat := scanner.Text()

	// Read user input for sorting option
	var sortingOption string
	for {
		fmt.Println("Size sorting options:")
		fmt.Println("1. Descending")
		fmt.Println("2. Ascending")
		fmt.Print("Enter a sorting option: ")
		scanner.Scan()
		sortingOption = scanner.Text()
		if sortingOption == "1" || sortingOption == "2" {
			break
		} else {
			fmt.Println("Wrong option")
		}
	}

	// Determine if sorting should be in descending order
	descending := sortingOption == "1"

	// Call the listFilesAndFolders function to list files in the specified directory
	filesBySize := listFilesAndFolders(directory, fileFormat, descending)

	// Read user input for duplicate check option
	var dupeOption string
	for {
		fmt.Println("Check for duplicates?")
		scanner.Scan()
		dupeOption = scanner.Text()
		if dupeOption == "Yes" || dupeOption == "No" {
			break
		} else {
			fmt.Println("Wrong option")
		}
	}

	// Determine if duplicate check should be performed
	dupeCheck := dupeOption == "Yes"

	// Call the hashFiles function to hash files from listFilesAndFolders
	hashFiles(filesBySize, dupeCheck)
}
