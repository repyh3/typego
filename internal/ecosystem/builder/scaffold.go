package builder

import (
	"bytes"
	_ "embed" // Embed the template
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/main.tmpl
var mainTmplStr string

type MainTemplateData struct {
	NamedImports map[string]string // Path -> Name
	Shims        map[string]string
	Bridge       string
}

// ScaffoldMain generates the main.go file in the specified directory
func ScaffoldMain(dir string, namedImports map[string]string, shims map[string]string, bridge string) error {
	tmpl, err := template.New("main").Parse(mainTmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse main template: %w", err)
	}

	data := MainTemplateData{
		NamedImports: namedImports,
		Shims:        shims,
		Bridge:       bridge,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	target := filepath.Join(dir, "main.go")
	if err := os.WriteFile(target, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}
