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
	"[": LeftBracket,
	"]": RightBracket,

	"==": Equals,
	"!=": NotEquals,
	"<=": LessThanOrEquals,
	">=": GreaterThanOrEquals,
	"<":  LessThan,
	">":  GreaterThan,

	"=": Assign,
}
