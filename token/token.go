package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foobar"

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
	EQ       = "EQ"
	NOT_EQ   = "NOT_EQ"
	BANG     = "BANG"
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	WHILE    = "WHILE"
	IN       = "IN"
	CONTINUE = "CONTINUE"
	BREAK    = "BREAK"
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
	"yeet":     RETURN,
	"stalk":    FOR,
	"onRepeat": WHILE,
	"in":       IN,
	"pass":     CONTINUE,
	"bounce":   BREAK,
	"is":       EQ,
	"aint":     NOT_EQ,
	"nah":      BANG,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
