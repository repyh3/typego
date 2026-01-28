package core

import (
	"reflect"
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindMap_StringKeys(b *testing.B) {
	vm := sobek.New()
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m[string(rune(i))] = i
	}

	val := reflect.ValueOf(m)
	// We create a fresh visited map for each run effectively,
	// but to avoid allocation noise in benchmark we can reuse if it's not mutated for this input.
	// bindMap doesn't mutate visited for non-pointer/non-struct types that take address.
	visited := make(map[uintptr]sobek.Value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bindMap(vm, val, visited)
	}
}

func BenchmarkBindMap_IntKeys(b *testing.B) {
	vm := sobek.New()
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i
	}

	val := reflect.ValueOf(m)
	visited := make(map[uintptr]sobek.Value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bindMap(vm, val, visited)
	}
}

func TestBindMap(t *testing.T) {
	vm := sobek.New()

	tests := []struct {
		name string
		input interface{}
		check func(*testing.T, sobek.Value)
	}{
		{
			name: "String Keys",
			input: map[string]int{"a": 1, "b": 2},
			check: func(t *testing.T, v sobek.Value) {
				obj := v.ToObject(vm)
				if obj.Get("a").ToInteger() != 1 {
					t.Errorf("expected a=1")
				}
			},
		},
		{
			name: "Int Keys",
			input: map[int]int{10: 1, 20: 2},
			check: func(t *testing.T, v sobek.Value) {
				obj := v.ToObject(vm)
				if obj.Get("10").ToInteger() != 1 {
					t.Errorf("expected 10=1")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := bindValue(vm, reflect.ValueOf(tt.input), make(map[uintptr]sobek.Value))
			if err != nil {
				t.Fatalf("bindValue error: %v", err)
			}
			tt.check(t, val)
		})
	}
}
