package core

import (
	"reflect"
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindMapIntKeys(b *testing.B) {
	vm := sobek.New()
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i
	}

	val := reflect.ValueOf(m)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We need a fresh visited map for each iteration to avoid cache hits
		// masking the allocation cost we want to measure, although bindMap
		// creates new objects anyway.
		visited := make(map[uintptr]sobek.Value)
		_, err := bindMap(vm, val, visited)
		if err != nil {
			b.Fatal(err)
		}
	}
}
