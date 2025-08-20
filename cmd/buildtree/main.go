package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/neomen/buildtree/internal/builder"
	"github.com/neomen/buildtree/internal/parser"
)

var (
	version = "-dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	filePath := flag.String("input-file", "", "Path to file containing directory structure")
	helpFlag := flag.Bool("help", false, "Show help")
	maxDepth := flag.Int("max-depth", 20, "Maximum nesting depth allowed (0 = no limit)")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.StringVar(filePath, "i", "", "Alias for --input-file")
	flag.IntVar(maxDepth, "d", 20, "Alias for --max-depth")
	flag.BoolVar(versionFlag, "v", false, "Alias for --version")
	flag.BoolVar(helpFlag, "h", false, "Alias for --help")
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if *versionFlag {
		fmt.Printf("buildtree v%s\nCommit: %s\nBuilt: %s\n", version, commit, date)
		return
	}

	input := getInput(*filePath)

	// Parse the input structure
	root, err := parser.ParseInput(input)
	if err != nil {
		log.Fatalf("Error parsing input: %v", err)
	}

	// Build the file structure
	if err := builder.BuildTree(root, *maxDepth); err != nil {
		log.Fatalf("Error building tree: %v", err)
	}
}

func getInput(filePath string) string {
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		return string(content)
	}

	args := flag.Args()
	if len(args) < 1 {
		printHelp()
		log.Fatal("Error: No input structure provided")
	}
	return args[0]
}

func printHelp() {
	fmt.Println("Buildtree - Instant Directory Tree Builder")
	fmt.Println("Usage: buildtree [OPTIONS] \"DIRECTORY_STRUCTURE\"")
	fmt.Println("Options:")
	fmt.Println("  -i, --input-file FILE	Read structure from file")
	fmt.Println("  -d, --max-depth N	Maximum nesting depth allowed (0=unlimited, default:20)")
	fmt.Println("  -h, --help		Show this help message")
	fmt.Println("  -v, --version		Show version information")
	fmt.Println("\nExamples:")
	fmt.Println("  buildtree \"project/\n├── src/\n│   └── main.go\"")
	fmt.Println("  buildtree --input-file structure.txt")
	fmt.Println("\nStructure format:")
	fmt.Println("  myproject/")
	fmt.Println("  ├── dir1/")
	fmt.Println("  │   ├── file1.txt")
	fmt.Println("  │   └── subdir/")
	fmt.Println("  └── file2.txt")
}
