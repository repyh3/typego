package core

import (
	"github.com/grafana/sobek/ast"
)

// Visitor is an interface for AST visitors
type Visitor interface {
	Visit(node ast.Node) []TextEdit
}

var Visitors []Visitor

func RegisterVisitor(v Visitor) {
	Visitors = append(Visitors, v)
}

func WalkAndCollect(node ast.Node, source string) []TextEdit {
	var allEdits []TextEdit

	// 1. Visit current node with registered visitors
	for _, v := range Visitors {
		edits := v.Visit(node)
		allEdits = append(allEdits, edits...)
	}

	// 2. Recurse children
	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Body {
			allEdits = append(allEdits, WalkAndCollect(stmt, source)...)
		}
	case *ast.BlockStatement:
		for _, stmt := range n.List {
			allEdits = append(allEdits, WalkAndCollect(stmt, source)...)
		}
	case *ast.ExpressionStatement:
		allEdits = append(allEdits, WalkAndCollect(n.Expression, source)...)
	case *ast.FunctionLiteral:
		allEdits = append(allEdits, WalkAndCollect(n.Body, source)...)
	case *ast.ArrowFunctionLiteral:
		allEdits = append(allEdits, WalkAndCollect(n.Body, source)...)
	case *ast.AssignExpression:
		allEdits = append(allEdits, WalkAndCollect(n.Left, source)...)
		allEdits = append(allEdits, WalkAndCollect(n.Right, source)...)
	case *ast.BinaryExpression:
		allEdits = append(allEdits, WalkAndCollect(n.Left, source)...)
		allEdits = append(allEdits, WalkAndCollect(n.Right, source)...)
	case *ast.CallExpression:
		allEdits = append(allEdits, WalkAndCollect(n.Callee, source)...)
		for _, arg := range n.ArgumentList {
			allEdits = append(allEdits, WalkAndCollect(arg, source)...)
		}
	case *ast.IfStatement:
		allEdits = append(allEdits, WalkAndCollect(n.Test, source)...)
		allEdits = append(allEdits, WalkAndCollect(n.Consequent, source)...)
		if n.Alternate != nil {
			allEdits = append(allEdits, WalkAndCollect(n.Alternate, source)...)
		}
	case *ast.ReturnStatement:
		if n.Argument != nil {
			allEdits = append(allEdits, WalkAndCollect(n.Argument, source)...)
		}
	case *ast.FunctionDeclaration:
		allEdits = append(allEdits, WalkAndCollect(n.Function, source)...)
	case *ast.ExportDeclaration:
		if n.HoistableDeclaration != nil {
			if n.HoistableDeclaration.FunctionDeclaration != nil {
				allEdits = append(allEdits, WalkAndCollect(n.HoistableDeclaration.FunctionDeclaration, source)...)
			}
		}
		if n.Variable != nil {
			allEdits = append(allEdits, WalkAndCollect(n.Variable, source)...)
		}
		if n.LexicalDeclaration != nil {
			allEdits = append(allEdits, WalkAndCollect(n.LexicalDeclaration, source)...)
		}

	case *ast.VariableStatement:
		for _, item := range n.List {
			allEdits = append(allEdits, WalkAndCollect(item, source)...)
		}
	case *ast.LexicalDeclaration:
		for _, item := range n.List {
			allEdits = append(allEdits, WalkAndCollect(item, source)...)
		}
	case *ast.Binding:
		if n.Initializer != nil {
			allEdits = append(allEdits, WalkAndCollect(n.Initializer, source)...)
		}
	case *ast.ObjectLiteral:
		for _, prop := range n.Value {
			switch p := prop.(type) {
			case *ast.PropertyKeyed:
				allEdits = append(allEdits, WalkAndCollect(p.Value, source)...)
			}
		}
	case *ast.ArrayLiteral:
		for _, item := range n.Value {
			if item != nil {
				allEdits = append(allEdits, WalkAndCollect(item, source)...)
			}
		}
	case *ast.DotExpression:
		allEdits = append(allEdits, WalkAndCollect(n.Left, source)...)
	case *ast.NewExpression:
		allEdits = append(allEdits, WalkAndCollect(n.Callee, source)...)
		for _, arg := range n.ArgumentList {
			allEdits = append(allEdits, WalkAndCollect(arg, source)...)
		}
	}

	return allEdits
}
