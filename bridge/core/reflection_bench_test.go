package core

import (
	"testing"

	"github.com/grafana/sobek"
)

func BenchmarkBindStruct_Map(b *testing.B) {
	vm := sobek.New()
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m["key"+string(rune(i))] = i
	}
	s := struct {
		M map[string]int
	}{
		M: m,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BindStruct(vm, "s", s)
	}
}

func BenchmarkBindStruct_ByteSlice(b *testing.B) {
	vm := sobek.New()
	data := make([]byte, 1024*1024) // 1MB
	s := struct {
		Data []byte
	}{
		Data: data,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BindStruct(vm, "s", s)
	}
}

func BenchmarkBindStruct_IntSlice(b *testing.B) {
	vm := sobek.New()
	data := make([]int, 1000)
	s := struct {
		Data []int
	}{
		Data: data,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BindStruct(vm, "s", s)
	}
}
