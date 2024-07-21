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
	OPEN_PAREN
	CLOSE_PAREN

	COLON_EQUALS // :=

	LET
)

var keywords = map[string]TokenId{
	"let": LET,
}
