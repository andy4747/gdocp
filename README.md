# GDOCP - Go Documentation Parser

GDOCP (Go Documentation Parser) is a powerful CLI tool and library for generating Markdown documentation from Go source code. It offers both single-file and recursive parsing options, as well as a built-in HTTP server for browsing generated documentation.

## Features

- Parse single Go files or entire directory structures
- Generate Markdown documentation from Go source code
- Built-in HTTP server for browsing documentation
- Customizable output directory
- Syntax highlighting in the web interface

## Installation

To Download GDOCP, use the following command:

```bash
go get -u github.com/andy4747/gdocp
```

To Install

```bash
go install github.com/andy4747/gdocp@latest
```

## Usage

### CLI

GDOCP can be used as a command-line tool with the following options:

```bash
gdocp [options]
```

Options:

- `-input`: Input Go file to parse (required if not using -r or -http)
- `-output`: Output markdown file (default: output.md)
- `-r`: Recursively parse Go files in subdirectories
- `-http`: Start HTTP server on the specified address (e.g., :6060)

Examples:

```bash
# Parse a single file
gdocp -input main.go -output main.md

# Recursively parse all Go files in the current directory and subdirectories
gdocp -r

# Start the HTTP server for browsing documentation
gdocp -http :6060
```

### Library

You can also use GDOCP as a library in your Go projects:

```go
import (
    "github.com/andy4747/gdocp/internal/parser"
    "github.com/andy4747/gdocp/internal/markdown"
)

// Parse a single file
fileInfo, err := parser.ParseFile("path/to/file.go")
mdContent, err := markdown.GenerateMarkdown(fileInfo)

// Generate documentation recursively
docMap, err := parser.GenerateMarkdownRecursively(".")
```

## Web Interface

When using the `-http` flag, GDOCP starts a web server that allows you to browse the generated documentation. The interface includes:

- A list of all parsed documents
- Syntax-highlighted Markdown rendering
- Easy navigation between files

## Contributing

Contributions to GDOCP are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

GDOCP uses the following open-source libraries:

- [marked.js](https://marked.js.org/) for Markdown rendering
- [highlight.js](https://highlightjs.org/) for syntax highlighting

## Author

[andy4747]

---

For more information, please open an issue on GitHub.
