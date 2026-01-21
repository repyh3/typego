package polyfills

import (
	"os"
	"runtime"
	"strings"

	"github.com/dop251/goja"
)

// EnableProcess injects the Node.js `process` global
func EnableProcess(vm *goja.Runtime) {
	proc := vm.NewObject()

	env := vm.NewObject()
	whitelist := map[string]bool{
		"PATH":     true,
		"LANG":     true,
		"PWD":      true,
		"HOSTNAME": true,
		"USER":     true,
	}

	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			upperKey := strings.ToUpper(key)
			// Allow whitelisted vars or anything prefixed with TYPEGO_
			if whitelist[upperKey] || strings.HasPrefix(upperKey, "TYPEGO_") {
				_ = env.Set(key, parts[1])
			}
		}
	}

	// Force color
	_ = env.Set("FORCE_COLOR", "1")
	_ = proc.Set("env", env)

	// process.platform
	_ = proc.Set("platform", runtime.GOOS)

	// process.cwd()
	_ = proc.Set("cwd", func(call goja.FunctionCall) goja.Value {
		wd, _ := os.Getwd()
		return vm.ToValue(wd)
	})

	// process.argv
	_ = proc.Set("argv", os.Args)

	// process.version
	_ = proc.Set("version", runtime.Version())

	_ = vm.Set("process", proc)
}
