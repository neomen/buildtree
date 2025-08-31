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
- **Smart filtering**: Include only relevant files by extension with `--filter`
- **Size control**: Limit file size with `--max-size` (default: 100kb)
- **Ignore patterns**: Skip directories like `.git` and `node_modules` by default
- **Hidden files**: Control inclusion of hidden files with `--include-hidden`

### Usage

```bash
buildtree -s OUTPUT_FILE [DIRECTORY] [OPTIONS]
```

This example will collect the contents of all files with the *.go, *.txt, *.md extension and combine them into a single file with relative paths.
```bash
buildtree -s export-structure.txt ./project-dir -f go,txt,md -I vendor
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
