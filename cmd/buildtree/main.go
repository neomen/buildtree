package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/neomen/buildtree/internal/builder"
	"github.com/neomen/buildtree/internal/parser"
)

var (
	version = "-dev"
	commit  = "none"
	date    = "unknown"
)

// Добавим интерфейсы для зависимостей, чтобы можно было мокировать их в тестах
type parserInterface interface {
	ParseInput(input string) (*parser.Node, error)
}

type builderInterface interface {
	BuildTree(root *parser.Node, maxDepth int) error
}

// Реальные реализации
type realParser struct{}
type realBuilder struct{}

func (r *realParser) ParseInput(input string) (*parser.Node, error) {
	return parser.ParseInput(input)
}

func (r *realBuilder) BuildTree(root *parser.Node, maxDepth int) error {
	return builder.BuildTree(root, maxDepth)
}

// Вынесем основную логику в отдельную функцию для тестирования
func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, p parserInterface, b builderInterface) int {
	// Создаем новый набор флагов для каждого вызова
	flags := flag.NewFlagSet("buildtree", flag.ContinueOnError)
	flags.SetOutput(stderr) // Устанавливаем вывод ошибок флагов в stderr

	filePath := flags.String("input-file", "", "Path to file containing directory structure")
	helpFlag := flags.Bool("help", false, "Show help")
	maxDepth := flags.Int("max-depth", 20, "Maximum nesting depth allowed (0 = no limit)")
	versionFlag := flags.Bool("version", false, "Show version information")
	flags.StringVar(filePath, "i", "", "Alias for --input-file")
	flags.IntVar(maxDepth, "d", 20, "Alias for --max-depth")
	flags.BoolVar(versionFlag, "v", false, "Alias for --version")
	flags.BoolVar(helpFlag, "h", false, "Alias for --help")

	// Парсим аргументы
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if *helpFlag {
		printHelp(stdout)
		return 0
	}

	if *versionFlag {
		fmt.Fprintf(stdout, "buildtree v%s\nCommit: %s\nBuilt: %s\n", version, commit, date)
		return 0
	}

	input := getInput(*filePath, stdin, flags, stderr)
	if input == "" {
		return 1
	}

	// Parse the input structure
	root, err := p.ParseInput(input)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing input: %v\n", err)
		return 1
	}

	// Build the file structure
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
