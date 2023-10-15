package evaluator

type BuiltinFunction func(args ...Object) (Object, error)

func (builtin BuiltinFunction) Type() string   { return "builtin" }
func (builtin BuiltinFunction) String() string { return "builtin function" }

var builtins = map[string]BuiltinFunction{
	"len": func(args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, IncorrectArgumentCountError{Expected: 1, Got: len(args)}
		}

		str, ok := args[0].(String)
		if !ok {
			return Null{}, nil
		}

		return Integer(len(str)), nil
	},
}
