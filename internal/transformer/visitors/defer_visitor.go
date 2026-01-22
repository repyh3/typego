package visitors

import (
	"github.com/grafana/sobek/ast"
	"github.com/repyh/typego/internal/transformer/core"
)

type DeferVisitor struct{}

func (v *DeferVisitor) Visit(node ast.Node) []core.TextEdit {
	switch fn := node.(type) {
	case *ast.FunctionLiteral:
		return v.transformFunction(fn.Body)
	case *ast.ArrowFunctionLiteral:
		if block, ok := fn.Body.(*ast.BlockStatement); ok {
			return v.transformFunction(block)
		}
	}
	return nil
}

func (v *DeferVisitor) transformFunction(body *ast.BlockStatement) []core.TextEdit {
	hasDefer := false
	scanForDefer(body, &hasDefer)

	if !hasDefer {
		return nil
	}

	var edits []core.TextEdit

	// 1. Wrapper Start
	edits = append(edits, core.TextEdit{
		Offset:  int(body.LeftBrace),
		Length:  0,
		NewText: " return typego.scope(function(__defer) { ",
	})

	// 2. Wrapper End
	edits = append(edits, core.TextEdit{
		Offset:  int(body.RightBrace) - 1,
		Length:  0,
		NewText: " }); ",
	})

	// 3. Replace 'defer' calls
	edits = append(edits, renameDefers(body)...)

	return edits
}

func scanForDefer(node ast.Node, found *bool) {
	if *found {
		return
	}

	switch node.(type) {
	case *ast.FunctionLiteral, *ast.ArrowFunctionLiteral:
		return
	}

	if call, ok := node.(*ast.CallExpression); ok {
		if ident, ok := call.Callee.(*ast.Identifier); ok {
			if ident.Name == "defer" {
				*found = true
				return
			}
		}
	}

	// Dynamic traversal fallback is handled by the walker for general nodes,
	// but for scanForDefer we manually recurse into common structures for speed
	if b, ok := node.(*ast.BlockStatement); ok {
		for _, s := range b.List {
			scanForDefer(s, found)
		}
	}
	if es, ok := node.(*ast.ExpressionStatement); ok {
		scanForDefer(es.Expression, found)
	}
}

func renameDefers(root ast.Node) []core.TextEdit {
	var edits []core.TextEdit

	var walker func(node ast.Node)
	walker = func(node ast.Node) {
		switch node.(type) {
		case *ast.FunctionLiteral, *ast.ArrowFunctionLiteral:
			return
		}

		if call, ok := node.(*ast.CallExpression); ok {
			if ident, ok := call.Callee.(*ast.Identifier); ok && ident.Name == "defer" {
				edits = append(edits, core.TextEdit{
					Offset:  int(ident.Idx) - 1,
					Length:  5,
					NewText: "__defer",
				})
			}
		}

		if b, ok := node.(*ast.BlockStatement); ok {
			for _, s := range b.List {
				walker(s)
			}
		}
		if es, ok := node.(*ast.ExpressionStatement); ok {
			walker(es.Expression)
		}
	}

	walker(root)
	return edits
}
