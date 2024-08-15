package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andy4747/gdocp/internal/markdown"
	"github.com/andy4747/gdocp/internal/parser"
)

var docMap map[string]parser.MarkdownContent

// createOutputDir ensures that the output directory exists and returns its path.
func createOutputDir() (string, error) {
	outputDir := "gdocp_out"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}
	return outputDir, nil
}

// processFile parses a single Go file and generates markdown output.
func processFile(inputFile, outputFile string) error {
	fileInfo, err := parser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse file: %v", err)
	}

	mdContent, err := markdown.GenerateMarkdown(fileInfo)
	if err != nil {
		return fmt.Errorf("failed to generate markdown: %v", err)
	}

	if mdContent == "" {
		return fmt.Errorf("skipped generating markdown for %s: Author or File is empty", inputFile)
	}

	err = markdown.WriteToFile(mdContent, outputFile)
	if err != nil {
		return fmt.Errorf("failed to write markdown file: %v", err)
	}

	fmt.Printf("Markdown file generated: %s\n", outputFile)
	return nil
}

func main() {
	inputFile := flag.String("input", "", "Input Go file to parse")
	outputFile := flag.String("output", "output.md", "Output markdown file")
	recursive := flag.Bool("r", false, "Recursively parse Go files in subdirectories")
	httpAddr := flag.String("http", "", "Start HTTP server on the specified address (e.g., :6060)")
	flag.Parse()

	if *inputFile == "" && !*recursive && *httpAddr == "" {
		log.Fatal("Either an input file must be specified, the -r flag must be used for recursive processing, or the -http flag must be used to start the HTTP server")
	}

	outputDir, err := createOutputDir()
	if err != nil {
		log.Fatalf("Failed to set up output directory: %v", err)
	}

	if *recursive || *httpAddr != "" {
		docMap, err = parser.GenerateMarkdownRecursively(".")
		if err != nil {
			log.Fatalf("Failed to process directory recursively: %v", err)
		}

		// Remove entries with empty content
		for id, content := range docMap {
			if content.Content == "" {
				delete(docMap, id)
			}
		}

		if *recursive {
			for _, content := range docMap {
				outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(content.Path), ".go")+".md")
				if err := markdown.WriteToFile(content.Content, outputFilePath); err != nil {
					log.Printf("Failed to write markdown file %s: %v", outputFilePath, err)
				} else {
					fmt.Printf("Markdown file generated: %s\n", outputFilePath)
				}
			}
		}
	} else if *inputFile != "" {
		// Process single file
		outputFilePath := filepath.Join(outputDir, *outputFile)
		if err := processFile(*inputFile, outputFilePath); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	// Start HTTP server if -http flag is provided
	if *httpAddr != "" {
		startHTTPServer(*httpAddr)
	}
}

func startHTTPServer(addr string) {
	http.HandleFunc("/", handleListDocs)
	http.HandleFunc("/doc", handleViewDoc)

	log.Printf("Starting HTTP server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleListDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Document List</title>
		<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
	</head>
	<body>
		<h1>Document List</h1>
		<ul>
		{{range $path, $content := .}}
			<li><a href="/doc?path={{$path}}">{{$content.Path}}</a></li>
		{{end}}
		</ul>
	</body>
	</html>
	`
	t, err := template.New("doclist").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, docMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleViewDoc(w http.ResponseWriter, r *http.Request) {
	docPath := r.URL.Query().Get("path")

	content, ok := docMap[docPath]
	if !ok {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>{{.Path}}</title>
        <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"></script>
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/github-dark.min.css">
        <style>
            body {
                font-family: Arial, sans-serif;
                line-height: 1.6;
                color: #333;
                max-width: 800px;
                margin: 0 auto;
                padding: 20px;
            }
            pre {
                background-color: #f4f4f4;
                border: 1px solid #ddd;
                border-radius: 4px;
                padding: 10px;
                overflow-x: auto;
            }
            code {
                font-family: 'Courier New', Courier, monospace;
            }
        </style>
    </head>
    <body>
        <h1>{{.Path}}</h1>
        <div id="content"></div>
        <script>
            // Set up marked.js to use highlight.js for code highlighting
            marked.setOptions({
                highlight: function(code, lang) {
                    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
                    return hljs.highlight(code, { language }).value;
                },
                langPrefix: 'hljs language-'
            });

            // Parse and render the Markdown content
            document.getElementById('content').innerHTML = marked.parse({{.Content}});

            // Apply highlighting to all code blocks
            document.querySelectorAll('pre code').forEach((block) => {
                hljs.highlightBlock(block);
            });
        </script>
    </body>
    </html>
    `
	t, err := template.New("docview").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
