package main

import (
	"banek/bytecode"
	"banek/vm"
	"encoding/gob"
	"os"
	"testing"
)

func BenchmarkEmulator(b *testing.B) {
	file, err := os.Open("../../examples/test.bac")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	var executable bytecode.Executable
	err = gob.NewDecoder(file).Decode(&executable)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := vm.Execute(executable)
		if err != nil {
			b.Error(err)
		}
	}
}
