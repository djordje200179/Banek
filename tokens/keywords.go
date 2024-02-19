package tokens

var keywords = map[string]Type{
	"func":   Func,
	"return": Return,

	"let": Let,
	"mut": Mut,

	"if":   If,
	"else": Else,
	"then": Then,

	"while": While,
	"do":    Do,
	"for":   For,

	"true":      Bool,
	"false":     Bool,
	"undefined": Undefined,
}

func LookupIdent(ident string) Type {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}

	return Ident
}
