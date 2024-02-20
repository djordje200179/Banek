package builtins

import (
	"banek/runtime/objs"
	"fmt"
	"strings"
)

func builtinPrint(args []objs.Obj) (objs.Obj, error) {
	var sb strings.Builder

	for i, arg := range args {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(arg.String())
	}

	fmt.Print(sb.String())

	return objs.Obj{}, nil
}

func builtinPrintln(args []objs.Obj) (objs.Obj, error) {
	var sb strings.Builder

	for i, arg := range args {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(arg.String())
	}
	sb.WriteByte('\n')

	fmt.Print(sb.String())

	return objs.Obj{}, nil
}

func builtinRead(_ []objs.Obj) (objs.Obj, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		return objs.Obj{}, err
	}

	return objs.MakeString(input), nil
}

func builtinReadln(_ []objs.Obj) (objs.Obj, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return objs.Obj{}, err
	}

	return objs.MakeString(input), nil
}
