package builder_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/repyh/typego/internal/builder"
)

func TestShimTemplate_Validity(t *testing.T) {
	// Mock values for the template
	imports := `_ "github.com/repyh/typego/bridge/modules/fmt"`
	jsCode := `"console.log('hello')"`
	bindings := "// bindings"
	memLimit := 64 * 1024 * 1024

	// Generate the code
	code := fmt.Sprintf(builder.ShimTemplate, imports, jsCode, bindings, memLimit)

	// Verify it's valid Go code
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "main.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("Generated shim code is invalid Go: %v\nCode Preview:\n%s", err, truncate(code, 500))
	}

	// Verify constraints
	if !strings.Contains(code, "engine.NewEngine") {
		t.Error("Generated code missing engine initialization")
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
