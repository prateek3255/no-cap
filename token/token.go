package token

type TokenType string

const (
	ILLEGAL = "illegal"
	EOF     = "end of file"

	// Identifiers + literals
	IDENT  = "identifier" // add, foobar, x, y, ...
	INT    = "integer"    // 1343456
	FLOAT  = "float"      // 1.3456
	STRING = "string"     // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	EQ       = "is"
	NOT_EQ   = "aint"
	BANG     = "nah"
	FUNCTION = "function"
	LET      = "fr"
	TRUE     = "noCap"
	FALSE    = "cap"
	IF       = "vibe"
	ELSE     = "nvm"
	ELSE_IF  = "unless"
	RETURN   = "yeet"
	FOR      = "stalk"
	WHILE    = "onRepeat"
	IN       = "in"
	CONTINUE = "pass"
	BREAK    = "bounce"
	NULL     = "ghosted"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"cook":     FUNCTION,
	"fr":       LET,
	"cap":      FALSE,
	"noCap":    TRUE,
	"vibe":     IF,
	"nvm":      ELSE,
	"unless":   ELSE_IF,
	"yeet":     RETURN,
	"stalk":    FOR,
	"onRepeat": WHILE,
	"in":       IN,
	"pass":     CONTINUE,
	"bounce":   BREAK,
	"is":       EQ,
	"aint":     NOT_EQ,
	"nah":      BANG,
	"ghosted":  NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
