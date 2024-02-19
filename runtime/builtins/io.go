package builtins

import (
	"banek/runtime"
	"banek/runtime/primitives"
	"fmt"
	"strings"
)

func builtinPrint(args []runtime.Obj) (runtime.Obj, error) {
	var sb strings.Builder

	for i, arg := range args {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(arg.String())
	}

	fmt.Print(sb.String())

	return primitives.Undefined{}, nil
}

func builtinPrintln(args []runtime.Obj) (runtime.Obj, error) {
	var sb strings.Builder

	for i, arg := range args {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(arg.String())
	}
	sb.WriteByte('\n')

	fmt.Print(sb.String())

	return primitives.Undefined{}, nil
}

func builtinRead(_ []runtime.Obj) (runtime.Obj, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return nil, err
	}

	return primitives.String(input), nil
}

func builtinReadln(_ []runtime.Obj) (runtime.Obj, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return nil, err
	}

	return primitives.String(input), nil
}
