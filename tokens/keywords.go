package tokens

var keywords = map[string]TokenType{
	"func": Function,
	"let":  Let,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}

	return Identifier
}
