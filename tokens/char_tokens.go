package tokens

var CharTokens = map[string]TokenType{
	"+": Plus,
	"-": Minus,
	"*": Asterisk,
	"/": Slash,

	"!": Bang,

	",": Comma,
	";": SemiColon,

	"(": LeftParenthesis,
	")": RightParenthesis,
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
