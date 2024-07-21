package compiler

import (
	"fmt"
	"script/ast"
	"script/lexer"
	"script/vm"
	"strconv"
)

func Compile(bytecode *vm.Bytecode, program *ast.Program) error {
	for _, stmt := range program.Statements {
		if err := compileStmt(bytecode, stmt); err != nil {
			return err
		}
	}
	return nil
}

func compileStmt(bc *vm.Bytecode, stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.DeclareStmt:
		return compileDeclareStmt(bc, s)
	case *ast.AssignStmt:
		return compileAssignStmt(bc, s)
	case *ast.BlockStmt:
		return compileBlockStmt(bc, s)
	case *ast.ConditionalStmt:
		return compileIfStmt(bc, s)
	default:
		return ast.NewNodeError(stmt, fmt.Sprintf("unknown statement type %T", stmt))
	}
}

func compileIfStmt(bc *vm.Bytecode, s *ast.ConditionalStmt) error {
	if err := compileExpr(bc, s.Cond); err != nil {
		return err
	}

	buffer := make(vm.Bytecode, 0)

	if err := compileBlockStmt(&buffer, s.Block); err != nil {
		return err
	}

	bc.Instruction(vm.JUMP_F, bc.Len()+buffer.Len()+1)
	bc.AppendBytecode(buffer)
// bc.Instruction(vm.LABEL, "if#")
	return nil
}

func compileDeclareStmt(bc *vm.Bytecode, s *ast.DeclareStmt) error {
	if err := compileExpr(bc, s.Expr); err != nil {
		return err
	}
	bc.Instruction(vm.DECLARE, s.Ident.Symbol)
	return nil
}

func compileAssignStmt(bc *vm.Bytecode, s *ast.AssignStmt) error {
	if err := compileExpr(bc, s.Expr); err != nil {
		return err
	}
	bc.Instruction(vm.STORE, s.Ident.Symbol)
	return nil
}

func compileBlockStmt(bc *vm.Bytecode, s *ast.BlockStmt) error {
	bc.Instruction(vm.ENTER, nil)
	for _, stmt := range s.Statements {
		if err := compileStmt(bc, stmt); err != nil {
			return err
		}
	}
	bc.Instruction(vm.LEAVE, nil)
	return nil
}

func compileExpr(bc *vm.Bytecode, expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if err := compileBinaryExpr(bc, e); err != nil {
			return err
		}
	case *ast.Number:
		if err := compileNumber(bc, e); err != nil {
			return err
		}
	case *ast.Identifier:
		bc.Instruction(vm.LOAD, e.Symbol)
	default:
		return ast.NewNodeError(expr, fmt.Sprintf("unknown expression type %T", expr))
	}
	return nil
}

func compileNumber(bc *vm.Bytecode, e *ast.Number) error {
	f, err := strconv.ParseFloat(e.Value, 64)
	if err != nil {
		panic(err)
	}
	bc.Instruction(vm.PUSH, f)
	return nil
}

func compileBinaryExpr(bc *vm.Bytecode, e *ast.BinaryExpr) error {
	if err := compileExpr(bc, e.Right); err != nil {
		return err
	}
	if err := compileExpr(bc, e.Left); err != nil {
		return err
	}
	switch e.Operator {
	case lexer.PLUS:
		bc.Instruction(vm.ADD, nil)
	case lexer.MINUS:
		bc.Instruction(vm.SUB, nil)
	case lexer.ASTERISK:
		bc.Instruction(vm.MUL, nil)
	case lexer.SLASH:
		bc.Instruction(vm.DIV, nil)
	default:
		return ast.NewNodeError(e, "unknown operator in binary expression")
	}

	return nil
}
