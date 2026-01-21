package linker

import (
	"fmt"
	"strings"
)

// GenerateShim creates the Go code to bind the package to the VM.
// Binds both top-level functions and exported structs.
func GenerateShim(info *PackageInfo, variableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n\t// Bind %s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("\t%s := eng.VM.NewObject()\n", variableName))

	// Bind top-level functions
	for _, fn := range info.Exports {
		sb.WriteString(fmt.Sprintf("\t%s.Set(%q, %s.%s)\n", variableName, fn.Name, info.Name, fn.Name))
	}

	// Bind struct constructors (factory functions)
	for _, st := range info.Structs {
		// Skip unexported structs (defensive check)
		if len(st.Name) == 0 || st.Name[0] < 'A' || st.Name[0] > 'Z' {
			continue
		}
		sb.WriteString(fmt.Sprintf("\t// Struct: %s\n", st.Name))
		sb.WriteString(fmt.Sprintf("\t%s.Set(\"New%s\", func() interface{} { return &%s.%s{} })\n",
			variableName, st.Name, info.Name, st.Name))
	}

	sb.WriteString(fmt.Sprintf("\teng.VM.Set(%q, %s)\n", "_go_hyper_"+info.Name, variableName))
	return sb.String()
}

// GenerateTypes creates the TypeScript definition with JSDoc.
func GenerateTypes(info *PackageInfo) string {
	var sb strings.Builder

	// Build known structs map for type resolution
	knownStructs := make(map[string]bool)
	for _, st := range info.Structs {
		knownStructs[st.Name] = true
	}

	sb.WriteString(fmt.Sprintf("// MODULE: go:%s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("declare module \"go:%s\" {\n", info.ImportPath))

	// Collect imports
	imports := make(map[string]map[string]bool) // pkgPath -> typeName -> bool
	for _, st := range info.Structs {
		for _, field := range st.Fields {
			if field.ImportPath != "" && field.ImportPath != info.ImportPath {
				if imports[field.ImportPath] == nil {
					imports[field.ImportPath] = make(map[string]bool)
				}
				typeName := field.Type
				if idx := strings.LastIndex(typeName, "."); idx != -1 {
					typeName = typeName[idx+1:]
				}
				typeName = strings.TrimPrefix(typeName, "*")
				typeName = strings.ReplaceAll(typeName, "[]", "")
				imports[field.ImportPath][typeName] = true
			}
		}
	}

	// Generate import statements
	for pkgPath, types := range imports {
		var typeList []string
		for t := range types {
			typeList = append(typeList, t)
		}
		sb.WriteString(fmt.Sprintf("\timport { %s } from \"go:%s\";\n", strings.Join(typeList, ", "), pkgPath))
	}
	if len(imports) > 0 {
		sb.WriteString("\n")
	}

	for _, st := range info.Structs {
		sb.WriteString(generateStructInterfaceWithContext(st, knownStructs))
	}

	for _, fn := range info.Exports {
		sb.WriteString(generateFunctionDeclWithContext(fn, knownStructs))
	}

	sb.WriteString("}\n")
	sb.WriteString(fmt.Sprintf("// END: go:%s\n", info.ImportPath))
	return sb.String()
}

func generateStructInterfaceWithContext(st ExportedStruct, knownStructs map[string]bool) string {
	var sb strings.Builder

	if st.Doc != "" {
		sb.WriteString("\t/**\n")
		for _, line := range strings.Split(st.Doc, "\n") {
			sb.WriteString(fmt.Sprintf("\t * %s\n", line))
		}
		sb.WriteString("\t */\n")
	}

	interfaceName := st.Name
	if len(st.TypeParams) > 0 {
		interfaceName += "<" + strings.Join(st.TypeParams, ", ") + ">"
	}

	// Handle embedded types with extends clause
	if len(st.Embeds) > 0 {
		// Filter to only include known structs from this package
		var validEmbeds []string
		for _, embed := range st.Embeds {
			if knownStructs == nil || knownStructs[embed] {
				validEmbeds = append(validEmbeds, embed)
			}
		}
		if len(validEmbeds) > 0 {
			sb.WriteString(fmt.Sprintf("\texport interface %s extends %s {\n", interfaceName, strings.Join(validEmbeds, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("\texport interface %s {\n", interfaceName))
		}
	} else {
		sb.WriteString(fmt.Sprintf("\texport interface %s {\n", interfaceName))
	}

	for _, field := range st.Fields {
		tsType := GoToTSTypeWithStructs(field.Type, knownStructs)
		if field.ImportPath != "" {
			if idx := strings.LastIndex(tsType, "."); idx != -1 {
				tsType = tsType[idx+1:]
			}
		}
		sb.WriteString(fmt.Sprintf("\t\t%s: %s;\n", field.Name, tsType))
	}

	for _, method := range st.Methods {
		sb.WriteString(generateMethodSignatureWithContext(method, knownStructs))
	}

	sb.WriteString("\t}\n\n")
	return sb.String()
}

func generateMethodSignatureWithContext(m MethodInfo, knownStructs map[string]bool) string {
	var args []string
	for i, arg := range m.Args {
		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}

		tsType := GoToTSTypeWithStructs(arg.Type, knownStructs)
		if strings.HasPrefix(arg.Type, "...") {
			name = "..." + name
		}

		args = append(args, fmt.Sprintf("%s: %s", name, tsType))
	}

	retType := "void"
	if len(m.Returns) > 0 {
		if len(m.Returns) == 1 {
			retType = GoToTSTypeWithStructs(m.Returns[0], knownStructs)
		} else {
			var types []string
			for _, r := range m.Returns {
				types = append(types, GoToTSTypeWithStructs(r, knownStructs))
			}
			retType = "[" + strings.Join(types, ", ") + "]"
		}
	}

	var sb strings.Builder
	if m.Doc != "" {
		sb.WriteString("\t\t/**\n")
		for _, line := range strings.Split(m.Doc, "\n") {
			sb.WriteString(fmt.Sprintf("\t\t * %s\n", line))
		}
		sb.WriteString("\t\t */\n")
	}
	sb.WriteString(fmt.Sprintf("\t\t%s(%s): %s;\n", m.Name, strings.Join(args, ", "), retType))
	return sb.String()
}

func generateFunctionDeclWithContext(fn ExportedFunc, knownStructs map[string]bool) string {
	var sb strings.Builder

	if fn.Doc != "" {
		sb.WriteString("\t/**\n")
		for _, line := range strings.Split(fn.Doc, "\n") {
			sb.WriteString(fmt.Sprintf("\t * %s\n", line))
		}
		sb.WriteString("\t */\n")
	}

	var args []string
	for i, arg := range fn.Args {
		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}

		tsType := GoToTSTypeWithStructs(arg.Type, knownStructs)
		if strings.HasPrefix(arg.Type, "...") {
			name = "..." + name
		}

		args = append(args, fmt.Sprintf("%s: %s", name, tsType))
	}

	retType := "void"

	returns := make([]string, len(fn.Ret))
	copy(returns, fn.Ret)

	if len(returns) > 0 && returns[len(returns)-1] == "error" {
		returns = returns[:len(returns)-1]
	}

	if len(returns) > 0 {
		if len(returns) == 1 {
			retType = GoToTSTypeWithStructs(returns[0], knownStructs)
		} else {
			var types []string
			for _, r := range returns {
				types = append(types, GoToTSTypeWithStructs(r, knownStructs))
			}
			retType = "[" + strings.Join(types, ", ") + "]"
		}
	}

	sb.WriteString(fmt.Sprintf("\texport function %s(%s): %s;\n", fn.Name, strings.Join(args, ", "), retType))
	return sb.String()
}

// GenerateTSShim creates the TypeScript source for the virtual module entry point.
// This is used by the compiler to resolve imports like "go:github.com/..."
func GenerateTSShim(info *PackageInfo) string {
	var sb strings.Builder
	for _, fn := range info.Exports {
		sb.WriteString(fmt.Sprintf("export const %s = (globalThis as any)._go_hyper_%s.%s;\n", fn.Name, info.Name, fn.Name))
	}
	// Add struct factory functions if any
	for _, st := range info.Structs {
		if len(st.Name) > 0 && st.Name[0] >= 'A' && st.Name[0] <= 'Z' {
			sb.WriteString(fmt.Sprintf("export const New%s = (globalThis as any)._go_hyper_%s.New%s;\n", st.Name, info.Name, st.Name))
		}
	}
	return sb.String()
}
