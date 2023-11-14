package builtins

import (
	"banek/runtime/objs"
	"banek/runtime/types"
	"fmt"
	"strings"
)

func builtinPrint(args []types.Obj) (types.Obj, error) {
	var sb strings.Builder

	for _, arg := range args {
		sb.WriteString(arg.String())
	}

	fmt.Print(sb.String())

	return objs.Undefined{}, nil
}

func builtinPrintln(args []types.Obj) (types.Obj, error) {
	var sb strings.Builder

	for _, arg := range args {
		sb.WriteString(arg.String())
	}

	fmt.Println(sb.String())

	return objs.Undefined{}, nil
}

func builtinRead(_ []types.Obj) (types.Obj, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return nil, err
	}

	return objs.Str(input), nil
}

func builtinReadln(_ []types.Obj) (types.Obj, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return nil, err
	}

	return objs.Str(input), nil
}
