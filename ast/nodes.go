package ast

import (
	"fmt"
	"script"
	"script/lexer"
)

// Node is an interface that represents a node in the AST.
type Node interface {
	fmt.Stringer
	Tok() lexer.Token
}

// Expr is an interface that represents an expression in the AST.
type Expr interface {
	Node
	expr()
}

// Stmt is an interface that represents a statement in the AST.
type Stmt interface {
	Node
	stmt()
}

type Program struct {
	Statements []Stmt
}

func (p *Program) Tok() lexer.Token {
	if len(p.Statements) == 0 {
		return lexer.Token{Id: lexer.EOF, Pos: -1}
	}
	return p.Statements[0].Tok()
}

func (p *Program) String() string {
	return script.Stringify(p)
}

type Identifier struct {
	Symbol string
	tok    lexer.Token
}

func (i *Identifier) Tok() lexer.Token {
	return i.tok
}

func (i *Identifier) String() string {
	return script.Stringify(i)
}

func (i *Identifier) expr() {}

type Number struct {
	Value string
	tok   lexer.Token
}

func (n *Number) Tok() lexer.Token {
	return n.tok
}

func (n *Number) String() string {
	return script.Stringify(n)
}

func (n *Number) expr() {}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.TokenId
	Right    Expr
}

func (b *BinaryExpr) Tok() lexer.Token {
	return b.Left.Tok()
}

func (b *BinaryExpr) String() string {
	return script.Stringify(b)
}

func (b *BinaryExpr) expr() {}

type UnaryExpr struct {
	Operator lexer.TokenId
	Expr     Expr
}

func (u *UnaryExpr) Tok() lexer.Token {
	return u.Expr.Tok()
}

func (u *UnaryExpr) String() string {
	return script.Stringify(u)
}

func (u *UnaryExpr) expr() {}

type FunctionExpr struct {
	Params []*Identifier
	Body   *BlockStmt
}

func (f *FunctionExpr) Tok() lexer.Token {
	return f.Params[0].Tok()
}

func (f *FunctionExpr) String() string {
	return script.Stringify(f)
}

func (f *FunctionExpr) expr() {}

//type ExprStmt struct {
//	Expr Expr
//}
//
//func (e *ExprStmt) Tok() lexer.Token {
//	return e.Expr.Tok()
//}
//
//func (e *ExprStmt) String() string {
//	return script.Stringify(e)
//}
//
//func (e *ExprStmt) stmt() {}

type DeclareStmt struct {
	Ident *Identifier
	Expr  Expr
}

func (d *DeclareStmt) Tok() lexer.Token {
	return d.Ident.Tok()
}

func (d *DeclareStmt) String() string {
	return script.Stringify(d)
}

func (d *DeclareStmt) stmt() {}

type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) Tok() lexer.Token {
	if len(b.Statements) == 0 {
		return lexer.Token{Id: lexer.EOF, Pos: -1}
	}
	return b.Statements[0].Tok()
}

func (b *BlockStmt) String() string {
	return script.Stringify(b)
}

func (b *BlockStmt) stmt() {}

type AssignStmt struct {
	Ident *Identifier
	Expr  Expr
}

func (a *AssignStmt) Tok() lexer.Token {
	return a.Ident.Tok()
}

func (a *AssignStmt) String() string {
	return script.Stringify(a)
}

func (a *AssignStmt) stmt() {}

type ConditionalStmt struct {
	Cond  Expr
	Block *BlockStmt
	Else  Stmt // TODO: else
}

func (i *ConditionalStmt) Tok() lexer.Token {
	return i.Cond.Tok()
}

func (i *ConditionalStmt) String() string {
	return script.Stringify(i)
}

func (i *ConditionalStmt) stmt() {}
