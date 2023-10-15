package tokens

var keywords = map[string]TokenType{
	"function": Function,
	"fn":       LambdaFunction,

	"let":   Let,
	"const": Const,

	"return": Return,

	"if":   If,
	"else": Else,
	"then": Then,

	"while": While,
	"true":  Boolean,
	"false": Boolean,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return Identifier
}
