package evaluator

type BuiltinFunction func(args ...Object) (Object, error)

func (builtin BuiltinFunction) Type() string   { return "builtin" }
func (builtin BuiltinFunction) String() string { return "builtin function" }

var builtins = map[string]BuiltinFunction{
	"len": func(args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, IncorrectArgumentCountError{Expected: 1, Got: len(args)}
		}

		switch arg := args[0].(type) {
		case String:
			return Integer(len(arg)), nil
		case Array:
			return Integer(len(arg)), nil
		default:
			return Null{}, nil
		}
	},
}
