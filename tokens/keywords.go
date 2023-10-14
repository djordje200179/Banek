package tokens

var keywords = map[string]TokenType{
	"function": Function,
	"let":      Let,
	"return":   Return,
	"if":       If,
	"else":     Else,
	"while":    While,
	"true":     Boolean,
	"false":    Boolean,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return Identifier
}
