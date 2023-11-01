package main

import (
	"banek/bytecode"
	"banek/vm"
	"encoding/gob"
	"os"
	"testing"
)

func BenchmarkEmulator(b *testing.B) {
	file, _ := os.Open("test.bac")
	defer file.Close()

	var executable bytecode.Executable
	_ = gob.NewDecoder(file).Decode(&executable)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = vm.Execute(executable)
	}
}
