package lexer

import (
	"fmt"
	"script"
)

func NewTokError(tok Token, message string) *script.PosError {
	e := &script.PosError{
		Pos:     tok.Pos,
		Message: fmt.Sprintf("(%s: %s): %s", tok.Id.String(), tok.Lexeme, message),
	}
	if script.PanicOnError {
		panic(e)
	}
	return e
}
