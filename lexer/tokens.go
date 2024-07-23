package lexer

import "fmt"

type Token struct {
	Pos    int
	Id     TokenId
	Lexeme string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s: %s)", t.Id.String(), t.Lexeme)
}

//go:generate stringer -type=TokenId
type TokenId int

const (
	INVALID TokenId = iota
	EOF
	LF

	IDENTIFIER
	NUMBER
	STRING
	CHAR

	PLUS
	MINUS
	ASTERISK
	SLASH
	EQUALS
	COLON
	COMMA

	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACE
	CLOSE_BRACE
	OPEN_BRACKET
	CLOSE_BRACKET

	COLON_EQUALS // :=

	EXCLAMATION

	EQUALS_EQUALS
	EXCLAMATION_EQUALS

	CIRCUMFLEX
	PIPE
	PIPE_PIPE
	AND
	AND_AND

	LESS_THAN
	GREATER_THAN
	LESS_THAN_EQUALS
	GREATER_THAN_EQUALS

	IF
	ELSE
	RETURN
	FOR
	CONTINUE
	BREAK
	FN
)

var keywords = map[string]TokenId{
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"for":      FOR,
	"continue": CONTINUE,
	"break":    BREAK,
	"fn":       FN,
}
