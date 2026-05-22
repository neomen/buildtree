package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/neomen/buildtree/internal/builder"
	"github.com/neomen/buildtree/internal/exporter"
	"github.com/neomen/buildtree/internal/parser"
	"github.com/neomen/buildtree/internal/utils"
)

var (
	version = "-dev"
	commit  = "none"
	date    = "unknown"
)

// Add interfaces for dependencies so that you can mock them in tests
type parserInterface interface {
	ParseInput(input string) (*parser.Node, error)
}

type builderInterface interface {
	BuildTree(root *parser.Node, maxDepth int) error
}

// Real implementations
type realParser struct{}
type realBuilder struct{}

func (r *realParser) ParseInput(input string) (*parser.Node, error) {
	return parser.ParseInput(input)
}

func (r *realBuilder) BuildTree(root *parser.Node, maxDepth int) error {
	return builder.BuildTree(root, maxDepth)
}

// Let's put the main logic in a separate function for testing
func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, p parserInterface, b builderInterface) int {
	// Parse flags manually from args (before -s handling) to handle flags that come after arguments
	var saveStructureValue string
	var filterValue string
	var ignoreDirsValue string
	var ignorePatternsValue string
	var maxSizeValue string
	var inputFilePath string
	var maxDepthValue string
	var includeHiddenValue string
	var minifyValue bool
	var helpFlagValue bool
	var versionFlagValue bool

	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Handle -s flag
		if arg == "-s" || arg == "--structure" {
			if i+1 < len(args) {
				saveStructureValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "-s=") || strings.HasPrefix(arg, "--structure=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				saveStructureValue = parts[1]
			}
		} else if arg == "-f" || arg == "--filter" {
			if i+1 < len(args) {
				filterValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "--filter=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				filterValue = parts[1]
			}
		} else if arg == "-I" || arg == "--ignore-dir" {
			if i+1 < len(args) {
				ignoreDirsValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "-I=") || strings.HasPrefix(arg, "--ignore-dir=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				ignoreDirsValue = parts[1]
			}
		} else if arg == "--ignore-pattern" || arg == "--ignore" {
			if i+1 < len(args) {
				ignorePatternsValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "--ignore-pattern=") || strings.HasPrefix(arg, "--ignore=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				ignorePatternsValue = parts[1]
			}
		} else if arg == "--max-size" {
			if i+1 < len(args) {
				maxSizeValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "--max-size=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				maxSizeValue = parts[1]
			}
		} else if arg == "--input-file" || arg == "-i" {
			if i+1 < len(args) {
				inputFilePath = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "--input-file=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				inputFilePath = parts[1]
			}
		} else if arg == "--max-depth" || arg == "-d" {
			if i+1 < len(args) {
				maxDepthValue = args[i+1]
				skipNext = true
			}
		} else if strings.HasPrefix(arg, "--max-depth=") || strings.HasPrefix(arg, "-d=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				maxDepthValue = parts[1]
			}
		} else if arg == "--include-hidden" || arg == "-H" {
			includeHiddenValue = "true"
		} else if strings.HasPrefix(arg, "--include-hidden=") || strings.HasPrefix(arg, "-H=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				includeHiddenValue = parts[1]
			}
		} else if arg == "--minify" {
			minifyValue = true
		} else if arg == "--help" || arg == "-h" {
			helpFlagValue = true
		} else if arg == "--version" || arg == "-v" {
			versionFlagValue = true
		}
	}

	// Processing of standard flags
	if helpFlagValue {
		printHelp(stdout)
		return 0
	}

	if versionFlagValue {
		fmt.Fprintf(stdout, "buildtree v%s\nCommit: %s\nBuilt: %s\n", version, commit, date)
		return 0
	}

	// FLAG HANDLING -s
	if saveStructureValue != "" {
		// Checking for conflicts
		if inputFilePath != "" {
			fmt.Fprintln(stderr, "Error: -s/--structure cannot be used with --input-file")
			return 1
		}

		// Defining the root directory
		rootDir := "."
		// Find the first non-flag argument after -s value
		skipNext := false
		for _, arg := range args {
			if skipNext {
				skipNext = false
				continue
			}
			if strings.HasPrefix(arg, "-") {
				// Skip the value if this is -s/--structure
				if arg == "-s" || arg == "--structure" || strings.HasPrefix(arg, "-s=") || strings.HasPrefix(arg, "--structure=") {
					skipNext = true
				}
				continue
			}
			rootDir = arg
			break
		}

		// Processing filters - default to text file extensions if not specified
		var filters []string
		if filterValue != "" {
			filters = strings.Split(filterValue, ",")
			for i := range filters {
				filters[i] = strings.TrimSpace(filters[i])
			}
		} else {
			// Default text file extensions for export
			filters = []string{"go", "js", "ts", "jsx", "tsx", "py", "rb", "php", "java", "c", "cpp", "h", "hpp", "cs", "rs", "swift", "kt", "scala", "html", "css", "scss", "sass", "json", "yaml", "yml", "xml", "md", "txt", "csv", "toml", "ini", "cfg", "conf"}
		}

		// Processing ignored directories
		var ignoreList []string
		if ignoreDirsValue != "" {
			ignoreList = strings.Split(ignoreDirsValue, ",")
			for i := range ignoreList {
				ignoreList[i] = strings.TrimSpace(ignoreList[i])
			}
		}

		// Processing ignore patterns
		var ignorePatternsList []string
		if ignorePatternsValue != "" {
			ignorePatternsList = strings.Split(ignorePatternsValue, ",")
			for i := range ignorePatternsList {
				ignorePatternsList[i] = strings.TrimSpace(ignorePatternsList[i])
			}
		}

		// Default max size if empty (use 100kb as default)
		maxSizeBytes, err := utils.ParseSize("100kb")
		if maxSizeValue != "" {
			maxSizeBytes, err = utils.ParseSize(maxSizeValue)
			if err != nil {
				fmt.Fprintf(stderr, "Error parsing max-size: %v\n", err)
				return 1
			}
		}

		// Check include hidden
		includeHidden := includeHiddenValue != "" && includeHiddenValue != "false"

		// Check minify
		minify := minifyValue

		// Export the structure
		if err := exporter.ExportStructure(saveStructureValue, rootDir, filters, ignoreList, ignorePatternsList, maxSizeBytes, includeHidden, minify); err != nil {
			fmt.Fprintf(stderr, "Error exporting structure: %v\n", err)
			return 1
		}
		return 0
	}

	// Normal processing (creating a structure)
	input := getInput(inputFilePath, stdin, args, stderr)
	if input == "" {
		return 1
	}

	// Parse max depth
	maxDepth := 20
	if maxDepthValue != "" {
		fmt.Sscanf(maxDepthValue, "%d", &maxDepth)
	}

	// Parsing the structure
	root, err := p.ParseInput(input)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing input: %v\n", err)
		return 1
	}

	// Building a tree
	if err := b.BuildTree(root, maxDepth); err != nil {
		fmt.Fprintf(stderr, "Error building tree: %v\n", err)
		return 1
	}

	return 0
}

func getInput(filePath string, stdin io.Reader, args []string, stderr io.Writer) string {
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(stderr, "Error reading file: %v\n", err)
			return ""
		}
		return string(content)
	}

	// Find the first non-flag argument
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return arg
	}

	printHelp(stderr)
	fmt.Fprintln(stderr, "Error: No input structure provided")
	return ""
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "Buildtree - Instant Directory Tree Builder")
	fmt.Fprintln(w, "Usage: buildtree [OPTIONS] \"DIRECTORY_STRUCTURE\"")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -i, --input-file FILE	Read structure from file")
	fmt.Fprintln(w, "  -d, --max-depth N	Maximum nesting depth allowed (0=unlimited, default:20)")
	fmt.Fprintln(w, "  -h, --help		Show this help message")
	fmt.Fprintln(w, "  -v, --version		Show version information")
	fmt.Fprintln(w, "  -s, --structure FILE	Save current directory structure with file contents to FILE")
	fmt.Fprintln(w, "  -f, --filter EXTS	Comma-separated list of file extensions or glob patterns (default: common text formats)")
	fmt.Fprintln(w, "  -I, --ignore-dir DIRS	Comma-separated list of directories to ignore (in addition to .git and node_modules)")
	fmt.Fprintln(w, "  --ignore PATTERNS	Comma-separated glob patterns for files to ignore (e.g., *_test.go,*.log)")
	fmt.Fprintln(w, "  --max-size SIZE	Maximum file size to include (default: 100kb, e.g., 100kb, 1mb)")
	fmt.Fprintln(w, "  -H, --include-hidden	Include hidden files and directories (starting with .)")
	fmt.Fprintln(w, "  --minify		Strip whitespace and empty lines to reduce token usage for LLM")
	fmt.Fprintln(w, "\nExamples:")
	fmt.Fprintln(w, "  buildtree \"project/\n├── src/\n│   └── main.go\"")
	fmt.Fprintln(w, "  buildtree --input-file structure.txt")
	fmt.Fprintln(w, "\nStructure format:")
	fmt.Fprintln(w, "  myproject/")
	fmt.Fprintln(w, "  ├── dir1/")
	fmt.Fprintln(w, "  │   ├── file1.txt")
	fmt.Fprintln(w, "  │   └── subdir/")
	fmt.Fprintln(w, "  └── file2.txt")
}

func main() {
	p := &realParser{}
	b := &realBuilder{}
	exitCode := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, p, b)
	os.Exit(exitCode)
}
