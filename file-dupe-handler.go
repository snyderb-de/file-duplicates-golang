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
	"strings"
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
		// This is the ascending order sorting
		sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
	}

	// Print files grouped by size
	for _, size := range sizes {
		fmt.Println(size, "bytes")
		for _, path := range filesBySize[size] {
			fmt.Println(path)
		}
		fmt.Println() // new line for grouping
	}

	return filesBySize
}

// Hash files from ListFilesAndFolders using MD5
func hashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			// Update err if no previous error has been encountered
			err = cerr
		}
	}()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), err
}

func hashFiles(filesBySize map[int64][]string, dupeCheck bool, descending bool) {
	if dupeCheck {
		sizes := make([]int64, 0, len(filesBySize))
		for size := range filesBySize {
			sizes = append(sizes, size)
		}

		if descending {
			sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })
		} else {
			sort.Slice(sizes, func(i, j int) bool { return sizes[i] < sizes[j] })
		}

		globalFileCounter := 1 // Initialize a global counter for continuous numbering

		for _, size := range sizes {
			files := filesBySize[size]
			if len(files) > 1 {
				filesByHash := make(map[string][]string)

				for _, file := range files {
					hash, err := hashFile(file)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}

					filesByHash[hash] = append(filesByHash[hash], file)
				}

				fmt.Printf("%d bytes\n", size) // Print size once for each group
				for hash, groupedFiles := range filesByHash {
					if len(groupedFiles) > 1 {
						fmt.Printf("Hash: %s\n", hash)
						for _, file := range groupedFiles {
							fmt.Printf("%d. %s\n", globalFileCounter, file)
							globalFileCounter++ // global counter
						}
					}
				}
				fmt.Println() // Add newline after each size group
			}
		}
	}
}

func main() {
	// Declare normalizedInput outside the for loop
	var normalizedInput string

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
		fmt.Println("Check for duplicates? (Yes/No)")
		scanner.Scan()
		dupeOption = scanner.Text()
		normalizedInput = strings.ToLower(dupeOption) // Normalize input to lowercase
		if normalizedInput == "yes" || normalizedInput == "no" {
			break
		} else {
			fmt.Println("Wrong option, please enter Yes or No")
		}
	}

	// Determine if duplicate check should be performed
	dupeCheck := normalizedInput == "yes"

	// Call the hashFiles function to hash files from listFilesAndFolders
	if dupeCheck {
		hashFiles(filesBySize, dupeCheck, descending)
	}
}
