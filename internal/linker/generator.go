package linker

import (
	"fmt"
	"strings"
)

func GenerateShim(info *PackageInfo, variableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n\t// Bind %s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("\t%s := eng.VM.NewObject()\n", variableName))

	for _, fn := range info.Exports {
		sb.WriteString(fmt.Sprintf("\t%s.Set(%q, %s.%s)\n", variableName, fn.Name, info.Name, fn.Name))
	}

	for _, st := range info.Structs {
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

func GenerateTypes(info *PackageInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// MODULE: go:%s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("declare module \"go:%s\" {\n", info.ImportPath))

	imports := make(map[string]map[string]bool)
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
		sb.WriteString(generateStructInterface(st))
	}

	for _, fn := range info.Exports {
		sb.WriteString(generateFunctionDecl(fn))
	}

	sb.WriteString("}\n")
	sb.WriteString(fmt.Sprintf("// END: go:%s\n", info.ImportPath))
	return sb.String()
}

func generateStructInterface(st ExportedStruct) string {
	var sb strings.Builder

	if st.Doc != "" {
		sb.WriteString("\t/**\n")
		for _, line := range strings.Split(st.Doc, "\n") {
			sb.WriteString(fmt.Sprintf("\t * %s\n", line))
		}
		sb.WriteString("\t */\n")
	}

	// Format interface name with generics: export interface Box<T> {
	interfaceName := st.Name
	if len(st.TypeParams) > 0 {
		interfaceName += "<" + strings.Join(st.TypeParams, ", ") + ">"
	}
	sb.WriteString(fmt.Sprintf("\texport interface %s {\n", interfaceName))

	for _, field := range st.Fields {
		tsType := field.TSType
		if field.ImportPath != "" {
			if idx := strings.LastIndex(tsType, "."); idx != -1 {
				tsType = tsType[idx+1:]
			}
		}
		sb.WriteString(fmt.Sprintf("\t\t%s: %s;\n", field.Name, tsType))
	}

	for _, method := range st.Methods {
		sb.WriteString(generateMethodSignature(method))
	}

	sb.WriteString("\t}\n\n")
	return sb.String()
}

func generateMethodSignature(m MethodInfo) string {
	var args []string
	for i, arg := range m.Args {
		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}

		tsType := GoToTSType(arg.Type)
		if strings.HasPrefix(arg.Type, "...") {
			name = "..." + name
		}

		args = append(args, fmt.Sprintf("%s: %s", name, tsType))
	}

	retType := "void"
	if len(m.Returns) > 0 {
		if len(m.Returns) == 1 {
			retType = GoToTSType(m.Returns[0])
		} else {
			var types []string
			for _, r := range m.Returns {
				types = append(types, GoToTSType(r))
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

func generateFunctionDecl(fn ExportedFunc) string {
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

		tsType := GoToTSType(arg.Type)
		if strings.HasPrefix(arg.Type, "...") {
			name = "..." + name
		}

		args = append(args, fmt.Sprintf("%s: %s", name, tsType))
	}

	retType := "void"

	// Goja runtime behavior:
	// - If (T, error), returns T or throws exception
	// - If (T1, T2, error), returns [T1, T2] or throws exception
	// - If method (via reflection.go), returns [T, Error] always

	// Create a copy of returns to manipulate
	returns := make([]string, len(fn.Ret))
	copy(returns, fn.Ret)

	// If last return is error, drop it (Goja throws it)
	if len(returns) > 0 && returns[len(returns)-1] == "error" {
		returns = returns[:len(returns)-1]
	}

	if len(returns) > 0 {
		if len(returns) == 1 {
			retType = GoToTSType(returns[0])
		} else {
			var types []string
			for _, r := range returns {
				types = append(types, GoToTSType(r))
			}
			retType = "[" + strings.Join(types, ", ") + "]"
		}
	}

	sb.WriteString(fmt.Sprintf("\texport function %s(%s): %s;\n", fn.Name, strings.Join(args, ", "), retType))
	return sb.String()
}
