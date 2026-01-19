package linker

import (
	"fmt"
	"strings"
)

// GenerateShim creates the Go code to bind the package to the VM.
func GenerateShim(info *PackageInfo, variableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n\t// Bind %s\n", info.ImportPath))
	sb.WriteString(fmt.Sprintf("\t%s := eng.VM.NewObject()\n", variableName))

	for _, fn := range info.Exports {
		sb.WriteString(fmt.Sprintf("\t%s.Set(%q, %s.%s)\n", variableName, fn.Name, info.Name, fn.Name))
	}

	sb.WriteString(fmt.Sprintf("\teng.VM.Set(%q, %s)\n", "_go_hyper_"+info.Name, variableName))
	return sb.String()
}

// GenerateTypes creates the TypeScript definition with JSDoc.
func GenerateTypes(info *PackageInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// MODULE: go:%s\n", info.Name))
	sb.WriteString(fmt.Sprintf("declare module \"go:%s\" {\n", info.ImportPath))

	for _, st := range info.Structs {
		sb.WriteString(generateStructInterface(st))
	}

	for _, fn := range info.Exports {
		sb.WriteString(generateFunctionDecl(fn))
	}

	sb.WriteString("}\n")
	sb.WriteString(fmt.Sprintf("// END: go:%s\n", info.Name))
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

	sb.WriteString(fmt.Sprintf("\texport interface %s {\n", st.Name))

	for _, field := range st.Fields {
		sb.WriteString(fmt.Sprintf("\t\t%s: %s;\n", field.Name, field.TSType))
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
		args = append(args, fmt.Sprintf("%s: %s", name, GoToTSType(arg.Type)))
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
		args = append(args, fmt.Sprintf("%s: %s", name, GoToTSType(arg.Type)))
	}

	sb.WriteString(fmt.Sprintf("\texport function %s(%s): unknown;\n", fn.Name, strings.Join(args, ", ")))
	return sb.String()
}
