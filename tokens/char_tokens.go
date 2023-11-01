package tokens

var CharTokens = map[string]TokenType{
	"+": Plus,
	"-": Minus,
	"*": Asterisk,
	"/": Slash,
	"%": Modulo,
	"^": Caret,

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

	"=":  Assign,
	"+=": PlusAssign,
	"-=": MinusAssign,
	"*=": AsteriskAssign,
	"/=": SlashAssign,
	"%=": ModuloAssign,
	"^=": CaretAssign,
}
