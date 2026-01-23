package core

import (
	"testing"

	"github.com/grafana/sobek"
)

func TestBindStruct_MapKeys(t *testing.T) {
	vm := sobek.New()
	data := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	err := BindStruct(vm, "testMap", data)
	if err != nil {
		t.Fatalf("BindStruct failed: %v", err)
	}

	val, err := vm.RunString(`testMap["1"]`)
	if err != nil {
		t.Fatalf("RunString failed: %v", err)
	}

	if val.String() != "one" {
		t.Errorf("Expected 'one', got %v", val)
	}

	val2, err := vm.RunString(`testMap["2"]`)
	if err != nil {
		t.Fatalf("RunString failed: %v", err)
	}
	if val2.String() != "two" {
		t.Errorf("Expected 'two', got %v", val2)
	}
}

func BenchmarkBindMap(b *testing.B) {
	vm := sobek.New()
	data := make(map[int]int)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "benchMap", data)
	}
}

func BenchmarkBindSlice(b *testing.B) {
	vm := sobek.New()
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BindStruct(vm, "benchSlice", data)
	}
}

func TestConsoleLog(t *testing.T) {
	vm := sobek.New()
	RegisterConsole(vm)

	// Just ensure it doesn't panic
	_, err := vm.RunString(`console.log("hello", "world", 123)`)
	if err != nil {
		t.Errorf("console.log failed: %v", err)
	}
}
