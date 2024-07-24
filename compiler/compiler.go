package compiler

import (
	"fmt"
	"script/ast"
	"script/lexer"
	"script/vm"
	"strconv"
)

type compiler struct {
	bytecode *vm.Bytecode
	// functions Stores the index of the function in the bytecode
	functions map[string]int
}

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
		return compileBlockStmt(bc, s, true)
	case *ast.ConditionalStmt:
		return compileConditionalStmt(bc, s)
	case *ast.ReturnStmt:
		return compileReturnStmt(bc, s)
	case *ast.ArrayAssignStmt:
		return compileArrayAssignStmt(bc, s)
	case *ast.ForStmt:
		return compileForStmt(bc, s)
	case *ast.ContinueStmt:
		return compileContinueStmt(bc, s)
	case *ast.BreakStmt:
		return compileBreakStmt(bc, s)
	default:
		return ast.NewNodeError(stmt, fmt.Sprintf("unknown statement type %T", stmt))
	}
}

func compileContinueStmt(bc *vm.Bytecode, s *ast.ContinueStmt) error {
	bc.Instruction(vm.JUMP_B, nil)
	return nil
}

func compileBreakStmt(bc *vm.Bytecode, s *ast.BreakStmt) error {
	bc.Instruction(vm.RET, nil)
	return nil
}

func compileConditionalStmt(bc *vm.Bytecode, s *ast.ConditionalStmt) error {
	if err := compileExpr(bc, s.Cond); err != nil {
		return err
	}

	jumpFalseIndex := bc.Len()
	bc.Instruction(vm.JUMP_F, -100)

	if err := compileBlockStmt(bc, s.Block, true); err != nil {
		return err
	}

	if s.Else != nil {
		jumpTrueIndex := bc.Len()
		bc.Instruction(vm.JUMP, -200)

		bc.SetArg(jumpFalseIndex, bc.Len())

		if err := compileStmt(bc, s.Else); err != nil {
			return err
		}
		bc.SetArg(jumpTrueIndex, bc.Len())
	} else {
		bc.SetArg(jumpFalseIndex, bc.Len())
	}

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
	if s.Ident != nil {
		bc.Instruction(vm.STORE, s.Ident.Symbol)
	} else {
		bc.Instruction(vm.POP, nil)
	}
	return nil
}

func compileBlockStmt(bc *vm.Bytecode, s *ast.BlockStmt, scope bool) error {
	if scope {
		bc.Instruction(vm.ENTER, nil)
	}

	for _, stmt := range s.Statements {
		if err := compileStmt(bc, stmt); err != nil {
			return err
		}
	}

	if scope {
		bc.Instruction(vm.LEAVE, nil)
	}
	return nil
}

func compileReturnStmt(bc *vm.Bytecode, s *ast.ReturnStmt) error {
	//// Push return expressions in reverse
	//for i := len(s.Returned) - 1; i >= 0; i-- {
	//	if err := compileExpr(bc, s.Returned[i]); err != nil {
	//		return err
	//	}
	//}
	//
	//// Push int with amount of arguments
	//bc.Instruction(vm.PUSH, len(s.Returned))

	// TODO Multiple return values
	if len(s.Returned) > 0 {

		if err := compileExpr(bc, s.Returned[0]); err != nil {
			return err
		}

	} else {
		bc.Instruction(vm.PUSH, nil)
	}
	bc.Instruction(vm.RET, nil)

	return nil
}

func compileForStmt(bc *vm.Bytecode, s *ast.ForStmt) error {
	frameIndex := bc.Len()
	bc.Instruction(vm.FRAME, -1)

	startIndex := bc.Len()
	bc.Instruction(vm.ENTER, nil)

	if s.Pre != nil {
		if err := compileStmt(bc, s.Pre); err != nil {
			return err
		}
	}

	innerFrameIndex := bc.Len()
	bc.Instruction(vm.FRAME, -1)

	if block, ok := s.Stmt.(*ast.BlockStmt); ok {
		if err := compileBlockStmt(bc, block, false); err != nil {
			return err
		}
	} else {
		if err := compileStmt(bc, block); err != nil {
			return err
		}
	}

	bc.Instruction(vm.LEAVE, nil)
	bc.Instruction(vm.JUMP, startIndex)

	bc.SetArg(frameIndex, bc.Len())

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
	case *ast.FunctionExpr:
		if err := compileFunctionExpr(bc, e); err != nil {
			return err
		}
	case *ast.CallExpr:
		if err := compileCallExpr(bc, e); err != nil {
			return err
		}
	case *ast.SubscriptExpr:
		if err := compileSubscriptExpr(bc, e); err != nil {
			return err
		}
	case *ast.ArrayExpr:
		if err := compileArrayExpr(bc, e); err != nil {
			return err
		}
	default:
		return ast.NewNodeError(expr, fmt.Sprintf("unknown expression type %T", expr))
	}
	return nil
}

func compileArrayAssignStmt(bc *vm.Bytecode, s *ast.ArrayAssignStmt) error {
	if err := compileExpr(bc, s.Expr); err != nil {
		return err
	}
	if err := compileExpr(bc, s.Ident); err != nil {
		return err
	}
	if err := compileExpr(bc, s.Index); err != nil {
		return err
	}
	bc.Instruction(vm.ARR_V, nil)
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
	default:
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

func compileFunctionExpr(bc *vm.Bytecode, e *ast.FunctionExpr) error {
	jumpIndex := bc.Len()
	bc.Instruction(vm.JUMP, -1)

	bc.Instruction(vm.ENTER, nil)

	argCountLabel := fmt.Sprintf("_argcount%d", bc.Len())

	if !e.IsVariadic {
		// Check if arg count matches exactly
		bc.Instruction(vm.PUSH, len(e.Params))
		bc.Instruction(vm.CMP, nil)
	} else {
		// Declare arg count
		bc.Instruction(vm.DECLARE, argCountLabel)

		bc.Instruction(vm.LOAD, argCountLabel)
		bc.Instruction(vm.PUSH, len(e.Params)-1) // Variadic args do accept an empty array.
		// Check if arg count matches or is greater
		bc.Instruction(vm.CMP_GTE, nil)
	}

	checkSuccessIndex := bc.Len()
	bc.Instruction(vm.JUMP_T, -1)

	// Return arg count error
	bc.Instruction(vm.PUSH, -1)
	compilePanic(bc, "arg count mismatch")

	bc.SetArg(checkSuccessIndex, bc.Len())

	for i, param := range e.Params {
		if e.IsVariadic && i >= len(e.Params)-1 {
			break
		}

		bc.Instruction(vm.DECLARE, param.Symbol)
	}

	if !e.IsVariadic {

	} else {
		index := len(e.Params) - 1

		// Calculate array size
		bc.Instruction(vm.LOAD, argCountLabel)
		bc.Instruction(vm.PUSH, len(e.Params)-1)
		bc.Instruction(vm.SUB, nil)

		bc.Instruction(vm.ARR_CR, nil)

		bc.Instruction(vm.DECLARE, e.Params[index].Symbol)
	}

	if err := compileBlockStmt(bc, e.Body, false); err != nil {
		return err
	}

	bc.Instruction(vm.LEAVE, nil)

	// TODO Will need to inject the return statement later. It IS required.
	//bc.Instruction(vm.PUSH, 0) // Returning no values
	//
	//bc.Instruction(vm.PUSH, nil) // Return nil
	//bc.Instruction(vm.RET, "NOTHING")

	bc.SetArg(jumpIndex, bc.Len())

	// Push index of function start. It is basically a pointer.

	function := vm.Func{
		Address: jumpIndex + 1,
	}

	bc.Instruction(vm.PUSH, function)

	return nil
}

func compilePanic(bc *vm.Bytecode, s string) {
	bc.Instruction(vm.PANIC, s)
}

func compileCallExpr(bc *vm.Bytecode, e *ast.CallExpr) error {
	// Push argument expressions in reverse to be declared in order in the call
	for i := len(e.Args) - 1; i >= 0; i-- {
		if err := compileExpr(bc, e.Args[i]); err != nil {
			return err
		}
	}

	// Arg count
	bc.Instruction(vm.PUSH, len(e.Args))

	if err := compileExpr(bc, e.Caller); err != nil {
		return err
	}

	frameReturnIndex := bc.Len()
	bc.Instruction(vm.FRAME, -1)

	bc.Instruction(vm.CALL, nil)
	bc.SetArg(frameReturnIndex, bc.Len())

	return nil
}

func compileArrayExpr(bc *vm.Bytecode, e *ast.ArrayExpr) error {
	// Push argument expressions in reverse to be in order
	for i := len(e.Elements) - 1; i >= 0; i-- {
		if err := compileExpr(bc, e.Elements[i]); err != nil {
			return err
		}
	}

	bc.Instruction(vm.PUSH, len(e.Elements))
	bc.Instruction(vm.ARR_CR, nil)
	return nil
}

func compileSubscriptExpr(bc *vm.Bytecode, e *ast.SubscriptExpr) error {
	if err := compileExpr(bc, e.Array); err != nil {
		return err
	}
	if err := compileExpr(bc, e.Index); err != nil {
		return err
	}
	bc.Instruction(vm.ARR_ID, nil)
	return nil
}
