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

func HandleUsageError() {
	flag.Usage()
	os.Exit(1)
}

func CheckError(err error) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprint(os.Stderr, err)
	os.Exit(2)
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
		HandleUsageError()
	}

	inputFileName := flag.Arg(0)
	inputFile, err := os.Open(inputFileName)
	CheckError(err)
	defer inputFile.Close()

	switch {
	case isInterpreted:
		err := Interpret(inputFile)
		CheckError(err)
	case isCompiled:
		outputFile, err := os.Create(outputFileName)
		CheckError(err)
		defer outputFile.Close()

		executable, err := Compile(inputFile)
		CheckError(err)

		RegisterGobTypes()
		encoder := gob.NewEncoder(outputFile)
		err = encoder.Encode(executable)
		CheckError(err)
	case isRun:
		var executable bytecode.Executable

		RegisterGobTypes()
		decoder := gob.NewDecoder(inputFile)
		err := decoder.Decode(&executable)
		CheckError(err)

		err = vm.Execute(executable)
		CheckError(err)
	default:
		HandleUsageError()
	}
}
