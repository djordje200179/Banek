package main

import (
	"banek/bytecode"
	"banek/exec/objects"
	"banek/vm"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func RegisterGobTypes() {
	gob.Register(objects.Array{})
	gob.Register(objects.Boolean(false))
	gob.Register(objects.BuiltinFunction{})
	gob.Register(objects.Integer(0))
	gob.Register(objects.String(""))
	gob.Register(objects.Undefined)
	gob.Register(objects.Unknown)
	gob.Register(&bytecode.Function{})
}

func main() {
	var isCompiled bool
	flag.BoolVar(&isCompiled, "c", false, "Compile the file")

	var isInterpreted bool
	flag.BoolVar(&isInterpreted, "i", false, "Interpret the file")

	var isRun bool
	flag.BoolVar(&isRun, "r", false, "Run the file")

	var outputFileName string
	flag.StringVar(&outputFileName, "o", "a.bac", "Output file name")

	flag.Parse()

	argsCount := flag.NArg()
	if argsCount == 0 || argsCount > 1 {
		flag.Usage()
		return
	}

	inputFileName := flag.Arg(0)
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer inputFile.Close()

	switch {
	case isInterpreted:
		err := Interpret(inputFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
	case isCompiled:
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		defer outputFile.Close()

		executable, err := Compile(inputFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}

		RegisterGobTypes()
		encoder := gob.NewEncoder(outputFile)
		err = encoder.Encode(executable)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
	case isRun:
		var executable bytecode.Executable

		RegisterGobTypes()
		decoder := gob.NewDecoder(inputFile)
		err := decoder.Decode(&executable)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}

		err = vm.Execute(executable)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
	default:
		flag.Usage()
		return
	}
}
