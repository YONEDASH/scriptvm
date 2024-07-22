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

	jumpIndex := bc.Len()
	bc.Instruction(vm.JUMP_F, -1)
	if err := compileBlockStmt(bc, s.Block); err != nil {
		return err
	}

	bc.SetArg(jumpIndex, bc.Len())

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
	case *ast.UnaryExpr:
		if err := compileUnaryExpr(bc, e); err != nil {
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
	// Check for operations that will short circuit
	switch e.Operator {
	case lexer.PIPE_PIPE:
		if err := compileExpr(bc, e.Left); err != nil {
			return err
		}
		jumpTrueIndex := bc.Len()
		bc.Instruction(vm.JUMP_T, -1)

		// If not true
		if err := compileExpr(bc, e.Right); err != nil {
			return err
		}
		secondTrueIndex := bc.Len()
		bc.Instruction(vm.JUMP_T, -1)

		// If false
		bc.Instruction(vm.PUSH, false) // <- will arrive here if false
		exitIndex := bc.Len()
		bc.Instruction(vm.JUMP, -1)

		// If true
		bc.SetArg(jumpTrueIndex, bc.Len())
		bc.SetArg(secondTrueIndex, bc.Len())
		bc.Instruction(vm.PUSH, true) // <- jump here if true
		bc.SetArg(exitIndex, bc.Len())

		return nil
	case lexer.AND_AND:
		if err := compileExpr(bc, e.Left); err != nil {
			return err
		}
		jumpFalseIndex := bc.Len()
		bc.Instruction(vm.JUMP_F, -1)

		// If not true
		if err := compileExpr(bc, e.Right); err != nil {
			return err
		}
		secondFalseIndex := bc.Len()
		bc.Instruction(vm.JUMP_F, -1)

		// If true
		bc.Instruction(vm.PUSH, true) // <- will arrive here if true
		exitIndex := bc.Len()
		bc.Instruction(vm.JUMP, -1)

		// If false
		bc.SetArg(jumpFalseIndex, bc.Len())
		bc.SetArg(secondFalseIndex, bc.Len())
		bc.Instruction(vm.PUSH, false) // <- jump here if false
		bc.SetArg(exitIndex, bc.Len())

		return nil
	}

	if err := compileExpr(bc, e.Left); err != nil {
		return err
	}
	if err := compileExpr(bc, e.Right); err != nil {
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
	case lexer.EQUALS_EQUALS:
		bc.Instruction(vm.CMP, nil)
	case lexer.EXCLAMATION_EQUALS:
		bc.Instruction(vm.CMP, nil)
		bc.Instruction(vm.NOT, nil)
	case lexer.LESS_THAN:
		bc.Instruction(vm.CMP_LT, nil)
	case lexer.GREATER_THAN:
		bc.Instruction(vm.CMP_GT, nil)
	case lexer.LESS_THAN_EQUALS:
		bc.Instruction(vm.CMP_LTE, nil)
	case lexer.GREATER_THAN_EQUALS:
		bc.Instruction(vm.CMP_GTE, nil)
	default:
		return ast.NewNodeError(e, "unknown operator in binary expression")
	}

	return nil
}

func compileUnaryExpr(bc *vm.Bytecode, e *ast.UnaryExpr) error {
	if err := compileExpr(bc, e.Expr); err != nil {
		return err
	}

	switch e.Operator {
	case lexer.EXCLAMATION:
		bc.Instruction(vm.NOT, nil)
	case lexer.MINUS:
		bc.Instruction(vm.NEG, nil)
	case lexer.PLUS:
		// Do nothing
	default:
		return ast.NewNodeError(e, "unknown operator in unary expression")
	}

	return nil
}
