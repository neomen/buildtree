### buildtree - Instant Directory Tree Builder

Buildtree is a lightning-fast CLI tool that creates directory structures from text-based tree diagrams. Perfect for developers, educators, and anyone who needs to quickly materialize directory examples provided by LLMs, documentation, or tutorials.

```text
custom_login/
├── custom_login.info.yml
├── custom_login.services.yml
├── src/
│   ├── Routing/
│   │   └── RouteSubscriber.php
│   └── Controller/
│       └── LoginController.php
```

→ Converts to actual file structure with one command!

---

## Features
- Instant creation of complex directory structures
- Zero dependencies - single binary
- Supports Windows, macOS, and Linux
- Handles Unicode tree characters (├─, └─, │)
- Automatically detects files vs directories

## Installation
```bash
go install github.com/neomen/buildtree/cmd/buildtree@latest
```

### Quick install (Linux/macOS)

```bash
# Using curl
curl -sSL https://github.com/neomen/buildtree/releases/latest/download/install.sh | sh

# Or using wget
wget -q -O - https://github.com/neomen/buildtree/releases/latest/download/install.sh | sh
```

### Manual installation

1. Download the latest release for your platform from the [releases page](https://github.com/neomen/buildtree/releases)
2. Extract the archive:
   ```bash
   tar xzf buildtree_<os>_<arch>.tar.gz
   ```
3. Move the binary to your PATH:
   ```bash
   # Linux/macOS
   sudo mv buildtree /usr/local/bin/

   # Windows
   # Move buildtree.exe to a directory in your PATH
   ```

## Usage

### Basic Example
```bash
buildtree "project/
├── src/
│   └── main.go
└── go.mod"
```

Creates:
```
project/
├── src
│   └── main.go
└── go.mod
```

### From Clipboard (Linux/macOS)
```bash
pbpaste | buildtree -
```

### From File
```bash
buildtree -i structure.txt
```

### Multi-level Project
```bash
buildtree "web-app/
├── public/
│   ├── index.html
│   └── assets/
│       ├── style.css
│       └── app.js
├── server/
│   └── main.py
└── README.md"
```

### Real-world LLM Example
```bash
buildtree "docker-project/
├── Dockerfile
├── .dockerignore
├── src/
│   └── app.py
├── requirements.txt
└── config/
    └── settings.yaml"
```

## Exporting Project Structure for LLM Context

Buildtree now includes a powerful feature that allows you to export your entire project structure with file contents to a single file, perfect for providing context to Large Language Models (LLMs).

### Key Features

- **Complete project context**: Export all files with their paths and contents in a structured format
- **Smart filtering**: Include only text files by default (go, js, py, md, json, yaml, xml, txt, etc.)
- **Size control**: Limit file size with `--max-size` (default: 100kb)
- **Ignore patterns**: Skip directories like `.git` and `node_modules` by default
- **Hidden files**: Control inclusion of hidden files with `--include-hidden`
- **File ignore patterns**: Ignore specific files using glob patterns with `--ignore`
- **Minification**: Reduce token usage with `--minify` (strips whitespace and empty lines)

### Usage

```bash
buildtree -s OUTPUT_FILE [DIRECTORY] [OPTIONS]
```

Options:
- `-f, --filter EXTS` - Comma-separated list of file extensions or glob patterns (default: common text formats like .go, .js, .txt, .md, .json, etc.)
- `-I, --ignore-dir DIRS` - Comma-separated list of directories to ignore (in addition to `.git` and `node_modules`)
- `--ignore PATTERNS` - Comma-separated glob patterns for files to ignore (e.g., `*_test.go,*.log`)
- `--max-size SIZE` - Maximum file size to include (default: 100kb, e.g., 100kb, 1mb)
- `-H, --include-hidden` - Include hidden files and directories (starting with `.`)
- `--minify` - Strip whitespace and empty lines to reduce token usage for LLM

Flags can be placed before or after the directory argument:
```bash
# Flag after directory
buildtree -s export.txt ./project-dir -f "go,md"

# Flag before directory
buildtree -s export.txt -f "go,md" ./project-dir
```

Examples:
```bash
# Export project structure (exports text files by default)
buildtree -s export-structure.txt ./project-dir -I vendor

# Export with glob patterns
buildtree -s all.txt . -f "*.go,*.md" --ignore="*_test.go,*.log"

# Export all files (don't filter by extension)
buildtree -s all.txt . -f "*"

# Export with multiple filters and ignore patterns
buildtree -s context.txt ./myapp -f "go,js,ts,md,json" --ignore="*_test.go,*.log,*.tmp" -I dist
```

### Advanced Filtering and Ignoring

**Glob pattern support**: Both `--filter` and `--ignore` support glob patterns like `*`, `?`, and `[...]`.

By default, only text file extensions are included: go, js, ts, py, rb, php, java, c, cpp, h, cs, rs, html, css, json, yaml, yml, xml, md, txt, csv, toml, ini, cfg, and more.

```bash
# Include only Go files and files starting with config
buildtree -s export.txt . -f "go,config.*"

# Ignore test files and log files
buildtree -s export.txt . --ignore="*_test.go,*.log"

# Multiple ignore patterns
buildtree -s export.txt . --ignore="*_test.go,*.log,*.tmp"

# Export all files without filtering by extension
buildtree -s all.txt . -f "*"

# Export with maximum file size limit
buildtree -s context.txt . --max-size 50kb

# Export including hidden files
buildtree -s all.txt . -f "*" -H

# Minify output to reduce token usage
buildtree -s context.txt . --minify
```

## Use Cases
- Quickly test LLM-generated file structures
- Create educational examples for documentation
- Reproduce project layouts during debugging
- Salvage directory structures from corrupted systems
- Automate project scaffolding in CI/CD pipelines

## Contribution
Contributions welcome! Please submit PRs for:
- Improved Unicode handling
- Windows clipboard integration
- Syntax extensions (file content hints)
- More robust error handling

---

**Stop manually creating directories** - let BuildTree materialize your file structures with magical speed!
