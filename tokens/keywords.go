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

	"true":      Boolean,
	"false":     Boolean,
	"undefined": Undefined,
}

func LookupIdentifier(identifier string) Type {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return Identifier
}
