package linker

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type PackageInfo struct {
	Name       string
	ImportPath string
	Exports    []ExportedFunc
	Structs    []ExportedStruct
}

type ExportedFunc struct {
	Name string
	Doc  string
	Args []ArgInfo
	Ret  []string
}

type ExportedStruct struct {
	Name        string
	PackagePath string // Source package import path
	TypeParams  []string
	Doc         string
	Fields      []FieldInfo
	Methods     []MethodInfo
	Embeds      []FieldInfo // Embedded types (for interface extension)
}

type FieldInfo struct {
	Name       string
	Type       string
	TSType     string
	IsNested   bool
	ImportPath string
}

type MethodInfo struct {
	Name           string
	Doc            string
	Receiver       string
	IsPointerRecv  bool
	Args           []ArgInfo
	Returns        []string
	HasCallbackArg bool
}

type ArgInfo struct {
	Name string
	Type string
}

func GoToTSType(goType string) string {
	return goToTSTypeWithContext(goType, nil)
}

func GoToTSTypeWithStructs(goType string, knownStructs map[string]bool) string {
	return goToTSTypeWithContext(goType, knownStructs)
}

func goToTSTypeWithContext(goType string, knownStructs map[string]bool) string {
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
	case "interface{}", "any":
		return "any"
	default:
		if strings.HasPrefix(goType, "[]") {
			return goToTSTypeWithContext(goType[2:], knownStructs) + "[]"
		}
		if strings.HasPrefix(goType, "map[") {
			return "Record<string, unknown>"
		}
		if strings.HasPrefix(goType, "func(") {
			return "(...args: unknown[]) => unknown"
		}
		if strings.HasPrefix(goType, "*") {
			return goToTSTypeWithContext(goType[1:], knownStructs)
		}
		// Handle variadic types
		if strings.HasPrefix(goType, "...") {
			baseType := goType[3:]
			return goToTSTypeWithContext(baseType, knownStructs) + "[]"
		}
		// Handle channels
		if strings.HasPrefix(goType, "chan ") || strings.HasPrefix(goType, "<-chan ") || strings.HasPrefix(goType, "chan<- ") {
			return "any"
		}
		// Handle anonymous structs
		if strings.HasPrefix(goType, "struct{") {
			return "any"
		}
		// Check if it's a known struct from this package
		if knownStructs != nil && knownStructs[goType] {
			return goType // Return as-is, it will reference the interface
		}
		// Return the type name (may be a struct from this package)
		return goType
	}
}

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
				parseTypeDecl(d, structMap, pkg.PkgPath, pkg.TypesInfo)
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

func parseTypeDecl(decl *ast.GenDecl, structMap map[string]*ExportedStruct, pkgPath string, info *types.Info) {
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
			Name:        ts.Name.Name,
			PackagePath: pkgPath,
			Doc:         strings.TrimSpace(decl.Doc.Text()),
		}

		if ts.TypeParams != nil {
			for _, param := range ts.TypeParams.List {
				for _, name := range param.Names {
					exported.TypeParams = append(exported.TypeParams, name.Name)
				}
			}
		}

		if st.Fields != nil {
			for _, field := range st.Fields.List {
				goType := types.ExprString(field.Type)

				// Embedded field (no name) - capture for interface extension
				if len(field.Names) == 0 {
					goType := types.ExprString(field.Type)
					cleanType := strings.TrimPrefix(goType, "*")

					var importPath string
					if info != nil {
						if t, ok := info.Types[field.Type]; ok {
							if named, ok := t.Type.(*types.Named); ok {
								if pkg := named.Obj().Pkg(); pkg != nil {
									importPath = pkg.Path()
								}
							} else if ptr, ok := t.Type.(*types.Pointer); ok {
								if named, ok := ptr.Elem().(*types.Named); ok {
									if pkg := named.Obj().Pkg(); pkg != nil {
										importPath = pkg.Path()
									}
								}
							}
						}
					}

					// Only add if it looks like a struct type (capitalized or qualified)
					if len(cleanType) > 0 {
						exported.Embeds = append(exported.Embeds, FieldInfo{
							Name:       cleanType,
							Type:       goType,
							ImportPath: importPath,
						})
					}
					continue
				}

				for _, name := range field.Names {
					if !ast.IsExported(name.Name) {
						continue
					}

					var importPath string
					if info != nil {
						if t, ok := info.Types[field.Type]; ok {
							if named, ok := t.Type.(*types.Named); ok {
								if pkg := named.Obj().Pkg(); pkg != nil {
									importPath = pkg.Path()
								}
							} else if ptr, ok := t.Type.(*types.Pointer); ok {
								if named, ok := ptr.Elem().(*types.Named); ok {
									if pkg := named.Obj().Pkg(); pkg != nil {
										importPath = pkg.Path()
									}
								}
							}
						}
					}

					exported.Fields = append(exported.Fields, FieldInfo{
						Name:       name.Name,
						Type:       goType,
						TSType:     GoToTSType(goType),
						IsNested:   isStructType(goType),
						ImportPath: importPath,
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

	var ret []string
	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			ret = append(ret, types.ExprString(result.Type))
		}
	}

	return ExportedFunc{
		Name: fn.Name.Name,
		Doc:  strings.TrimSpace(fn.Doc.Text()),
		Args: args,
		Ret:  ret,
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
