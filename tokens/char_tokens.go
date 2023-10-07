package tokens

var CharTokens = map[string]TokenType{
	"+": Plus,
	"-": Minus,
	"*": Multiply,
	"/": Divide,

	"!": Bang,

	",": Comma,
	";": SemiColon,

	"(": LeftParen,
	")": RightParen,
	"{": LeftBrace,
	"}": RightBrace,

	"==": Equals,
	"!=": NotEquals,
	"<=": LessThanOrEquals,
	">=": GreaterThanOrEquals,
	"<":  LessThan,
	">":  GreaterThan,

	"=": Assign,
}
