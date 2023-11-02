package tokens

var keywords = map[string]TokenType{
	"function": Function,
	"fn":       LambdaFunction,

	"let": Let,
	"mut": Mut,

	"return": Return,

	"if":   If,
	"else": Else,
	"then": Then,

	"while": While,
	"do":    Do,

	"true":      Boolean,
	"false":     Boolean,
	"undefined": Undefined,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return Identifier
}
