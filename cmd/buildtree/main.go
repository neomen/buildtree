package main

import (
	"flag"
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
	// First, we look for the -s or --structure flag in the arguments
	var saveStructureValue string
	var newArgs []string
	skipNext := false

	// We go through all the arguments to find the -s flag
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		// Checking the long shape of the flag
		if arg == "-s" || arg == "--structure" {
			if i+1 < len(args) {
				saveStructureValue = args[i+1]
				skipNext = true
				continue
			}
			// Checking the form with equality
		} else if strings.HasPrefix(arg, "-s=") || strings.HasPrefix(arg, "--structure=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				saveStructureValue = parts[1]
				continue
			}
		}

		// Saving arguments, except for the found flag
		newArgs = append(newArgs, arg)
	}

	// Creating a new set of flags
	flags := flag.NewFlagSet("buildtree", flag.ContinueOnError)
	flags.SetOutput(stderr)

	// Defining all flags
	filePath := flags.String("input-file", "", "Path to file containing directory structure")
	helpFlag := flags.Bool("help", false, "Show help")
	maxDepth := flags.Int("max-depth", 20, "Maximum nesting depth allowed (0 = no limit)")
	versionFlag := flags.Bool("version", false, "Show version information")

	// A flag for exporting the structure (needed for parsing other flags)
	saveStructure := flags.String("s", "", "Save current directory structure with file contents to specified file")

	// Additional flags for export
	structureFilter := flags.String("filter", "", "Comma-separated list of file extensions to include (e.g., go,js,txt)")
	ignoreDirs := flags.String("ignore-dir", "", "Comma-separated list of directories to ignore (in addition to .git and node_modules)")
	maxSize := flags.String("max-size", "100kb", "Maximum file size to include (e.g., 100kb, 1mb)")
	includeHidden := flags.Bool("include-hidden", false, "Include hidden files and directories (starting with .)")

	// Adding Aliases
	flags.StringVar(filePath, "i", "", "Alias for --input-file")
	flags.IntVar(maxDepth, "d", 20, "Alias for --max-depth")
	flags.BoolVar(versionFlag, "v", false, "Alias for --version")
	flags.BoolVar(helpFlag, "h", false, "Alias for --help")
	flags.StringVar(saveStructure, "structure", "", "Alias for --save-structure")
	flags.StringVar(structureFilter, "f", "", "Alias for --filter")
	flags.StringVar(ignoreDirs, "I", "", "Alias for --ignore-dir (uppercase i)")

	// Parse the remaining arguments
	if err := flags.Parse(newArgs); err != nil {
		return 1
	}

	// If we found the -s flag manually, we use its value
	if saveStructureValue != "" {
		*saveStructure = saveStructureValue
	}

	// Processing of standard flags
	if *helpFlag {
		printHelp(stdout)
		return 0
	}

	if *versionFlag {
		fmt.Fprintf(stdout, "buildtree v%s\nCommit: %s\nBuilt: %s\n", version, commit, date)
		return 0
	}

	// FLAG HANDLING -s
	if *saveStructure != "" {
		// Checking for conflicts
		if *filePath != "" {
			fmt.Fprintln(stderr, "Error: -s/--structure cannot be used with --input-file")
			return 1
		}

		// Defining the root directory
		rootDir := "."
		args := flags.Args()
		if len(args) > 0 {
			rootDir = args[0]
		}

		// Processing filters
		var filters []string
		if *structureFilter != "" {
			filters = strings.Split(*structureFilter, ",")
			for i := range filters {
				filters[i] = strings.TrimSpace(filters[i])
			}
		}

		// Processing ignored directories
		var ignoreList []string
		if *ignoreDirs != "" {
			ignoreList = strings.Split(*ignoreDirs, ",")
			for i := range ignoreList {
				ignoreList[i] = strings.TrimSpace(ignoreList[i])
			}
		}

		// Parse the maximum size
		maxSizeBytes, err := utils.ParseSize(*maxSize)
		if err != nil {
			fmt.Fprintf(stderr, "Error parsing max-size: %v\n", err)
			return 1
		}

		// Export the structure
		if err := exporter.ExportStructure(*saveStructure, rootDir, filters, ignoreList, maxSizeBytes, *includeHidden); err != nil {
			fmt.Fprintf(stderr, "Error exporting structure: %v\n", err)
			return 1
		}
		return 0
	}

	// Normal processing (creating a structure)
	input := getInput(*filePath, stdin, flags, stderr)
	if input == "" {
		return 1
	}

	// Parsing the structure
	root, err := p.ParseInput(input)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing input: %v\n", err)
		return 1
	}

	// Building a tree
	if err := b.BuildTree(root, *maxDepth); err != nil {
		fmt.Fprintf(stderr, "Error building tree: %v\n", err)
		return 1
	}

	return 0
}

func getInput(filePath string, stdin io.Reader, flags *flag.FlagSet, stderr io.Writer) string {
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(stderr, "Error reading file: %v\n", err)
			return ""
		}
		return string(content)
	}

	args := flags.Args()
	if len(args) < 1 {
		printHelp(stderr)
		fmt.Fprintln(stderr, "Error: No input structure provided")
		return ""
	}
	return args[0]
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
	fmt.Fprintln(w, "  -f, --filter EXTS	Comma-separated list of file extensions to include (e.g., go,js,txt)")
	fmt.Fprintln(w, "  -I, --ignore-dir DIRS	Comma-separated list of directories to ignore (in addition to .git and node_modules)")
	fmt.Fprintln(w, "  --max-size SIZE	Maximum file size to include (default: 100kb, e.g., 100kb, 1mb)")
	fmt.Fprintln(w, "  -H, --include-hidden	Include hidden files and directories (starting with .)")
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
