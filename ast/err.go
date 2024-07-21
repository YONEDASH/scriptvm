package ast

import (
	"script"
	"script/lexer"
)

func NewNodeError(node Node, message string) *script.PosError {
	return lexer.NewTokError(node.Tok(), message)
}
