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
go install github.com/neomen/buildtree@latest
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
buildtree -f structure.txt
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