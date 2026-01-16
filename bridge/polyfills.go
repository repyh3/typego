package bridge

import (
	"os"
	"runtime"
	"strings"

	"github.com/dop251/goja"
)

// EnableNodeCompat injects Node.js globals like process
func EnableNodeCompat(vm *goja.Runtime) {
	proc := vm.NewObject()

	// process.env
	env := vm.NewObject()
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			env.Set(parts[0], parts[1])
		}
	}
	// Force Node libraries (like chalk) to detect color support
	env.Set("FORCE_COLOR", "1")
	proc.Set("env", env)

	// process.platform
	proc.Set("platform", runtime.GOOS)

	// process.cwd()
	proc.Set("cwd", func(call goja.FunctionCall) goja.Value {
		wd, _ := os.Getwd()
		return vm.ToValue(wd)
	})

	// process.argv
	proc.Set("argv", os.Args)

	// process.version
	proc.Set("version", runtime.Version())

	vm.Set("process", proc)
}
