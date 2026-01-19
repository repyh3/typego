package linker

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// PackageInfo contains metadata about an inspected Go package.
type PackageInfo struct {
	Name       string
	ImportPath string
	Exports    []ExportedFunc
	Structs    []ExportedStruct
}

// ExportedFunc represents a top-level exported function.
type ExportedFunc struct {
	Name string
	Doc  string
	Args []ArgInfo
	Ret  []string
}

// ExportedStruct represents an exported struct type with its fields and methods.
type ExportedStruct struct {
	Name    string
	Doc     string
	Fields  []FieldInfo
	Methods []MethodInfo
}

// FieldInfo describes a struct field.
type FieldInfo struct {
	Name     string
	Type     string
	TSType   string
	IsNested bool
}

// MethodInfo describes a method attached to a struct.
type MethodInfo struct {
	Name           string
	Doc            string
	Receiver       string
	IsPointerRecv  bool
	Args           []ArgInfo
	Returns        []string
	HasCallbackArg bool
}

// ArgInfo describes a function/method argument.
type ArgInfo struct {
	Name string
	Type string
}

// GoToTSType maps Go types to TypeScript equivalents.
func GoToTSType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "error":
		return "Error | null"
	default:
		if strings.HasPrefix(goType, "[]") {
			return GoToTSType(goType[2:]) + "[]"
		}
		if strings.HasPrefix(goType, "map[") {
			return "Record<string, unknown>"
		}
		if strings.HasPrefix(goType, "func(") {
			return "(...args: unknown[]) => unknown"
		}
		if strings.HasPrefix(goType, "*") {
			return GoToTSType(goType[1:])
		}
		return goType
	}
}

// Inspect loads and analyzes a Go package for functions, structs, and methods.
func Inspect(importPath string, dir string) (*PackageInfo, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir: dir,
	}

	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %v", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package load errors")
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found")
	}

	pkg := pkgs[0]
	info := &PackageInfo{
		Name:       pkg.Name,
		ImportPath: pkg.PkgPath,
	}

	structMap := make(map[string]*ExportedStruct)

	for _, syntax := range pkg.Syntax {
		for _, decl := range syntax.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				parseTypeDecl(d, structMap)
			case *ast.FuncDecl:
				if d.Recv == nil {
					if ast.IsExported(d.Name.Name) {
						info.Exports = append(info.Exports, parseFunc(d))
					}
				} else {
					parseMethod(d, structMap)
				}
			}
		}
	}

	for _, s := range structMap {
		info.Structs = append(info.Structs, *s)
	}

	return info, nil
}

func parseTypeDecl(decl *ast.GenDecl, structMap map[string]*ExportedStruct) {
	for _, spec := range decl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok || !ast.IsExported(ts.Name.Name) {
			continue
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}

		exported := &ExportedStruct{
			Name: ts.Name.Name,
			Doc:  strings.TrimSpace(decl.Doc.Text()),
		}

		if st.Fields != nil {
			for _, field := range st.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				for _, name := range field.Names {
					if !ast.IsExported(name.Name) {
						continue
					}
					goType := types.ExprString(field.Type)
					exported.Fields = append(exported.Fields, FieldInfo{
						Name:     name.Name,
						Type:     goType,
						TSType:   GoToTSType(goType),
						IsNested: isStructType(goType),
					})
				}
			}
		}

		structMap[ts.Name.Name] = exported
	}
}

func parseFunc(fn *ast.FuncDecl) ExportedFunc {
	var args []ArgInfo
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			typeName := types.ExprString(param.Type)
			if len(param.Names) == 0 {
				args = append(args, ArgInfo{Name: "", Type: typeName})
			} else {
				for _, name := range param.Names {
					args = append(args, ArgInfo{Name: name.Name, Type: typeName})
				}
			}
		}
	}

	return ExportedFunc{
		Name: fn.Name.Name,
		Doc:  strings.TrimSpace(fn.Doc.Text()),
		Args: args,
	}
}

func parseMethod(fn *ast.FuncDecl, structMap map[string]*ExportedStruct) {
	if !ast.IsExported(fn.Name.Name) || fn.Recv == nil || len(fn.Recv.List) == 0 {
		return
	}

	recv := fn.Recv.List[0]
	recvType := types.ExprString(recv.Type)
	isPointer := strings.HasPrefix(recvType, "*")
	baseName := strings.TrimPrefix(recvType, "*")

	if !ast.IsExported(fn.Name.Name) {
		return
	}

	method := MethodInfo{
		Name:          fn.Name.Name,
		Doc:           strings.TrimSpace(fn.Doc.Text()),
		Receiver:      recvType,
		IsPointerRecv: isPointer,
	}

	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			typeName := types.ExprString(param.Type)
			if strings.HasPrefix(typeName, "func(") {
				method.HasCallbackArg = true
			}
			if len(param.Names) == 0 {
				method.Args = append(method.Args, ArgInfo{Name: "", Type: typeName})
			} else {
				for _, name := range param.Names {
					method.Args = append(method.Args, ArgInfo{Name: name.Name, Type: typeName})
				}
			}
		}
	}

	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			method.Returns = append(method.Returns, types.ExprString(result.Type))
		}
	}

	// Only attach method if the receiver type is already in structMap
	// (meaning it was identified as a struct type by parseTypeDecl)
	if existing, ok := structMap[baseName]; ok {
		existing.Methods = append(existing.Methods, method)
	}
	// Non-struct types (like "type Error string") are ignored
}

func isStructType(typeName string) bool {
	typeName = strings.TrimPrefix(typeName, "*")
	if len(typeName) == 0 {
		return false
	}
	first := typeName[0]
	return first >= 'A' && first <= 'Z'
}
