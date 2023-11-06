package tokens

var CharTokens = map[string]Type{
	"+": Plus,
	"-": Minus,
	"*": Asterisk,
	"/": Slash,
	"%": Modulo,
	"^": Caret,
	"!": Bang,

	"->": Arrow,
	"|":  VerticalBar,

	",": Comma,
	";": SemiColon,

	"(": LeftParen,
	")": RightParen,
	"{": LeftBrace,
	"}": RightBrace,
	"[": LeftBracket,
	"]": RightBracket,

	"==": Equals,
	"!=": NotEquals,
	"<=": LessEquals,
	">=": GreaterEquals,
	"<":  Less,
	">":  Greater,

	"=":  Assign,
	"+=": PlusAssign,
	"-=": MinusAssign,
	"*=": AsteriskAssign,
	"/=": SlashAssign,
	"%=": ModuloAssign,
	"^=": CaretAssign,
}
