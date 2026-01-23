package core

import (
	"fmt"
	"sort"

	"github.com/grafana/sobek/parser"
)

// Transform parses the source, applies visitors, and returns the modified source.
func Transform(filename, source string) (string, error) {
	prog, err := parser.ParseFile(nil, filename, source, 0)
	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	edits := WalkAndCollect(prog, source)

	// 3. Sort edits by offset to apply correctly in reverse
	sort.Slice(edits, func(i, j int) bool {
		if edits[i].Offset == edits[j].Offset {
			return edits[i].Length < edits[j].Length
		}
		return edits[i].Offset < edits[j].Offset
	})

	out := source
	for i := len(edits) - 1; i >= 0; i-- {
		edit := edits[i]
		out = out[:edit.Offset] + edit.NewText + out[edit.Offset+edit.Length:]
	}

	return out, nil
}

type TextEdit struct {
	Offset  int
	Length  int
	NewText string
}
