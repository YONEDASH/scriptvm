package ast

import (
	"errors"
	"fmt"
	"script/lexer"
)

func Parse(tokens []lexer.Token) (*Program, []error) {
	p := &parser{tokens: tokens, errors: make([]error, 0)}

	program := &Program{
		Statements: make([]Stmt, 0),
	}

	for {
		id := p.get(0).Id
		if id == lexer.EOF {
			break
		}

		if id == lexer.LF {
			p.index++
			continue
		}

		node, err := p.parseStmt()
		if err != nil {
			p.errors = append(p.errors, err)
			p.recover()
			continue
		}
		if node == nil {
			continue
		}

		if err = p.expectLF(); err != nil {
			p.errors = append(p.errors, err)
			p.recover()
			continue
		}

		program.Statements = append(program.Statements, node)
	}

	return program, p.errors
}

type parser struct {
	tokens []lexer.Token
	errors []error
	index  int
}

func (p *parser) get(offset int) lexer.Token {
	if p.index+offset >= len(p.tokens) {
		return lexer.Token{Id: lexer.EOF, Pos: -1}
	}
	return p.tokens[p.index+offset]
}

func (p *parser) consume() lexer.Token {
	t := p.get(0)
	p.index++
	return t
}

func (p *parser) done() bool {
	return p.get(0).Id == lexer.EOF
}

func (p *parser) recover() {
	for {
		id := p.get(0).Id
		if id == lexer.EOF || id == lexer.LF {
			break
		}
		p.index++
	}
}

func (p *parser) expect(id lexer.TokenId, msg string) (lexer.Token, error) {
	t := p.consume()
	if t.Id != id {
		return t, lexer.NewTokError(t, fmt.Sprintf("expected %s (%s), got %s", msg, id.String(), t.Id.String()))
	}
	return t, nil
}

func (p *parser) expectLF() error {
	switch p.get(0).Id {
	case lexer.LF:
		_, err := p.expect(lexer.LF, "LF")
		return err
	default:
		_, err := p.expect(lexer.EOF, "LF or EOF")
		return err
	}
}

func (p *parser) expectLFB() error {
	switch p.get(0).Id {
	case lexer.LF:
		_, err := p.expect(lexer.LF, "LF")
		return err
	case lexer.CLOSE_BRACE:
		_, err := p.expect(lexer.LF, "Close Brace")
		return err
	default:
		_, err := p.expect(lexer.EOF, "LF, close brace or EOF")
		return err
	}
}

func (p *parser) parseStmt() (Stmt, error) {
	switch p.get(0).Id {
	case lexer.LF:
		p.index++
		return nil, nil
	case lexer.OPEN_BRACE:
		return p.parseBlockStmt()
	case lexer.IF:
		return p.parseIfStmt()
	default:
	}

	switch p.get(1).Id {
	case lexer.COLON_EQUALS:
		return p.parseDeclareStmt()
	case lexer.EQUALS:
		return p.parseAssignStmt()
	default:
		return nil, lexer.NewTokError(p.get(0), "expected statement")
	}
}

func (p *parser) parseIfStmt() (Stmt, error) {
	if _, err := p.expect(lexer.IF, "expected if"); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "if statement"))
	}

	block, err := p.parseBlockStmt()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "if statement"))
	}

	return &ConditionalStmt{
		Cond:  expr,
		Block: block,
	}, nil
}

func (p *parser) parseDeclareStmt() (Stmt, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "expected identifier"))
	}

	if _, err = p.expect(lexer.COLON_EQUALS, "expected :="); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "expected expression"))
	}

	return &DeclareStmt{
		Ident: ident,
		Expr:  expr,
	}, nil
}

func (p *parser) parseAssignStmt() (Stmt, error) {
	ident, err := p.parseIdent()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "expected identifier"))
	}

	if _, err = p.expect(lexer.EQUALS, "expected ="); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, errors.Join(err, lexer.NewTokError(p.get(0), "expected expression"))
	}

	return &AssignStmt{
		Ident: ident,
		Expr:  expr,
	}, nil
}

func (p *parser) parseBlockStmt() (*BlockStmt, error) {
	if _, err := p.expect(lexer.OPEN_BRACE, "open brace"); err != nil {
		return nil, err
	}

	block := &BlockStmt{
		Statements: make([]Stmt, 0),
	}

	for {
		if p.get(0).Id == lexer.CLOSE_BRACE {
			break
		}

		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		if stmt == nil {
			continue
		}

		if err = p.expectLFB(); err != nil {
			return nil, err
		}

		block.Statements = append(block.Statements, stmt)
	}

	if _, err := p.expect(lexer.CLOSE_BRACE, "close brace"); err != nil {
		return nil, err
	}

	return block, nil
}

func (p *parser) parseExpr() (Expr, error) {
	return p.parseBinaryExprAdditive()
}

func (p *parser) parseBinaryExprAdditive() (Expr, error) {
	left, err := p.parseBinaryExprMultiplicative()
	if err != nil {
		return nil, err
	}

	for p.get(0).Id == lexer.PLUS || p.get(0).Id == lexer.MINUS {
		operator := p.consume()
		var right Expr
		if right, err = p.parseBinaryExprMultiplicative(); err != nil {
			return nil, err
		}

		left = &BinaryExpr{
			Left:     left,
			Operator: operator.Id,
			Right:    right,
		}
	}

	return left, nil
}

func (p *parser) parseBinaryExprMultiplicative() (Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.get(0).Id == lexer.ASTERISK || p.get(0).Id == lexer.SLASH {
		operator := p.consume()
		var right Expr
		if right, err = p.parsePrimary(); err != nil {
			return nil, err
		}

		left = &BinaryExpr{
			Left:     left,
			Operator: operator.Id,
			Right:    right,
		}
	}

	return left, nil
}

//
//func (p *parser) parseIdents() ([]*Identifier, error) {
//	if _, err := p.expect(lexer.OPEN_PAREN, "expected open paren for identifiers"); err != nil {
//		return nil, err
//	}
//
//	idents := make([]*Identifier, 0)
//	for !p.done() && p.get(0).Id != lexer.CLOSE_PAREN {
//		ident, err := p.parseIdent()
//		if err != nil {
//			return nil, errors.Join(err, lexer.NewTokError(p.get(0), "expected identifier"))
//		}
//		idents = append(idents, ident)
//	}
//
//	if _, err := p.expect(lexer.CLOSE_PAREN, "expected close paren for identifiers"); err != nil {
//		return nil, err
//	}
//
//	return idents, nil
//}

func (p *parser) parsePrimary() (Expr, error) {
	tk := p.get(0)

	// Ignore LF
	for ; tk.Id == lexer.LF; tk = p.get(0) {
		p.index++
	}

	switch tk.Id {
	case lexer.IDENTIFIER:
		return p.parseIdent()
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.OPEN_PAREN:
		//i := 0
		//for {
		//	id := p.get(i).Id
		//	if id == lexer.EOF || id == lexer.LF || id == lexer.CLOSE_PAREN {
		//		break
		//	}
		//	i++
		//}
		//if p.get(i).Id == lexer.CLOSE_PAREN && p.get(i+1).Id == lexer.OPEN_BRACE {
		//	idents, err := p.parseIdents()
		//	if err != nil {
		//		return nil, err
		//	}
		//	block, err := p.parseBlockStmt()
		//	return &FunctionExpr{Params: idents, Body: block}, err
		//}

		p.consume()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(lexer.CLOSE_PAREN, ")"); err != nil {
			return nil, err
		}
		return expr, nil
	default:
		p.index++
		return nil, lexer.NewTokError(tk, "expected primary expression")
	}
}

func (p *parser) parseIdent() (*Identifier, error) {
	t, err := p.expect(lexer.IDENTIFIER, "identifier")
	if err != nil {
		return nil, err
	}
	return &Identifier{Symbol: t.Lexeme, tok: t}, nil
}

func (p *parser) parseNumber() (*Number, error) {
	t, err := p.expect(lexer.NUMBER, "number")
	if err != nil {
		return nil, err
	}
	return &Number{Value: t.Lexeme, tok: t}, nil
}
