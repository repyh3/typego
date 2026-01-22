package visitors

import (
	"strconv"

	"github.com/grafana/sobek/ast"
	"github.com/repyh/typego/internal/transformer/core"
)

type IotaVisitor struct {
	counter int
}

func (v *IotaVisitor) Visit(node ast.Node) []core.TextEdit {
	var edits []core.TextEdit

	switch n := node.(type) {
	case *ast.VariableStatement:
		for _, decl := range n.List {
			edits = append(edits, v.checkDeclaration(decl)...)
		}
	case *ast.LexicalDeclaration:
		// LexicalDeclaration covers 'const' and 'let'
		// In JS, 'const' is a LexicalDeclaration with Token == CONST
		for _, decl := range n.List {
			edits = append(edits, v.checkDeclaration(decl)...)
		}
	case *ast.Program:
		// Reset counter at start of file
		v.counter = 0
	}

	return edits
}

func (v *IotaVisitor) checkDeclaration(decl *ast.Binding) []core.TextEdit {
	if decl.Initializer == nil {
		return nil
	}

	if ident, ok := decl.Initializer.(*ast.Identifier); ok && ident.Name == "iota" {
		edit := core.TextEdit{
			Offset:  int(ident.Idx) - 1,
			Length:  4,
			NewText: strconv.Itoa(v.counter),
		}
		v.counter++
		return []core.TextEdit{edit}
	}

	return nil
}
