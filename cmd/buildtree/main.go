package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Define command-line flags
	filePath := flag.String("f", "", "Path to file containing directory structure")
	helpFlag := flag.Bool("h", false, "Show help")
	helpFlagLong := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Show help if requested
	if *helpFlag || *helpFlagLong {
		printHelp()
		return
	}

	// Get input structure
	var input string
	if *filePath != "" {
		// Read from file
		content, err := os.ReadFile(*filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		input = string(content)
	} else {
		// Get from command-line arguments
		args := flag.Args()
		if len(args) < 1 {
			printHelp()
			log.Fatal("Error: No input structure provided")
		}
		input = args[0]
	}

	// Process the input
	processInput(input)
}

func printHelp() {
	fmt.Println("Buildtree - Instant Directory Tree Builder")
	fmt.Println("Usage: buildtree [OPTIONS] \"DIRECTORY_STRUCTURE\"")
	fmt.Println("Options:")
	fmt.Println("  -f FILE    Read structure from file")
	fmt.Println("  -h, --help Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  buildtree \"project/\n├── src/\n│   └── main.go\"")
	fmt.Println("  buildtree -f structure.txt")
	fmt.Println("\nStructure format:")
	fmt.Println("  myproject/")
	fmt.Println("  ├── dir1/")
	fmt.Println("  │   ├── file1.txt")
	fmt.Println("  │   └── subdir/")
	fmt.Println("  └── file2.txt")
	fmt.Printf("BuildTree v%s\n", version)
}

func processInput(input string) {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		log.Fatal("Input is empty")
	}

	// Process root directory
	rootDir := strings.TrimSpace(lines[0])
	// Remove comments from root directory
	if idx := strings.Index(rootDir, "#"); idx != -1 {
		rootDir = strings.TrimSpace(rootDir[:idx])
	}
	if strings.HasSuffix(rootDir, "/") {
		rootDir = rootDir[:len(rootDir)-1]
	}

	// Create root directory
	if err := os.Mkdir(rootDir, 0755); err != nil && !os.IsExist(err) {
		log.Fatal("Error creating root directory:", err)
	}

	// Stack to track parent directories for each level
	stack := []string{rootDir}

	for i, line := range lines[1:] {
		// Remove comments (everything after #)
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimRight(line, " ")
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Calculate indentation level
		depth := 0
		remaining := line
		for len(remaining) > 0 {
			r, size := utf8.DecodeRuneInString(remaining)
			if r == ' ' || r == '│' || r == '├' || r == '└' || r == '─' {
				depth++
				remaining = remaining[size:]
			} else {
				break
			}
		}

		// Calculate nesting level (4 characters = 1 level)
		level := depth / 4

		// Extract element name
		name := strings.TrimSpace(remaining)
		name = strings.TrimPrefix(name, "── ")
		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		// Validate level
		if level > len(stack) {
			log.Fatalf("Line %d: invalid nesting level %d (max %d)", i+2, level, len(stack))
		}

		// Determine parent directory
		parent := rootDir
		if level > 0 {
			parent = stack[level-1]
		}

		if strings.HasSuffix(name, "/") {
			// Create directory
			dirName := strings.TrimSuffix(name, "/")
			dirPath := filepath.Join(parent, dirName)
			if err := os.Mkdir(dirPath, 0755); err != nil && !os.IsExist(err) {
				log.Fatal("Error creating directory:", err)
			}

			// Update stack for this level
			if level < len(stack) {
				stack[level] = dirPath
			} else {
				stack = append(stack, dirPath)
			}
		} else {
			// Create file
			filePath := filepath.Join(parent, name)
			if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
				log.Fatal("Error creating file:", err)
			}
		}
	}
}
