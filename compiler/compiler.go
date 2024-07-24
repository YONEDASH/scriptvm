package compiler

import (
	"fmt"
	"script/ast"
	"script/lexer"
	"script/vm"
	"strconv"
)

type compiler struct {
	bc        *vm.Bytecode
	loopBegin stack[int]
	loopEnd   stack[int]
}

func Compile(bytecode *vm.Bytecode, program *ast.Program) error {
	c := &compiler{
		bc:        bytecode,
		loopBegin: make(stack[int], 4),
		loopEnd:   make(stack[int], 4),
	}

	for _, stmt := range program.Statements {
		if err := c.compileStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (out *compiler) compileStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.DeclareStmt:
		return out.compileDeclareStmt(s)
	case *ast.AssignStmt:
		return out.compileAssignStmt(s)
	case *ast.BlockStmt:
		return out.compileBlockStmt(s, true)
	case *ast.ConditionalStmt:
		return out.compileConditionalStmt(s)
	case *ast.ReturnStmt:
		return out.compileReturnStmt(s)
	case *ast.ArrayAssignStmt:
		return out.compileArrayAssignStmt(s)
	case *ast.ForStmt:
		return out.compileForStmt(s)
	case *ast.ContinueStmt:
		return out.compileContinueStmt(s)
	case *ast.BreakStmt:
		return out.compileBreakStmt(s)
	default:
		return ast.NewNodeError(stmt, fmt.Sprintf("unknown statement type %T", stmt))
	}
}

func (out *compiler) compileContinueStmt(s *ast.ContinueStmt) error {
	out.bc.Instruction(vm.JUMP, out.loopBegin.top())
	return nil
}

func (out *compiler) compileBreakStmt(s *ast.BreakStmt) error {
	out.bc.Instruction(vm.JUMP, out.loopEnd.top())
	return nil
}

func (out *compiler) compileConditionalStmt(s *ast.ConditionalStmt) error {
	if err := out.compileExpr(s.Cond); err != nil {
		return err
	}

	jumpFalseIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP_F, -100)

	if err := out.compileBlockStmt(s.Block, true); err != nil {
		return err
	}

	if s.Else != nil {
		jumpTrueIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP, -200)

		out.bc.SetArg(jumpFalseIndex, out.bc.Len())

		if err := out.compileStmt(s.Else); err != nil {
			return err
		}
		out.bc.SetArg(jumpTrueIndex, out.bc.Len())
	} else {
		out.bc.SetArg(jumpFalseIndex, out.bc.Len())
	}

	return nil
}

func (out *compiler) compileDeclareStmt(s *ast.DeclareStmt) error {
	if err := out.compileExpr(s.Expr); err != nil {
		return err
	}
	out.bc.Instruction(vm.DECLARE, s.Ident.Symbol)
	return nil
}

func (out *compiler) compileAssignStmt(s *ast.AssignStmt) error {
	if err := out.compileExpr(s.Expr); err != nil {
		return err
	}
	if s.Ident != nil {
		out.bc.Instruction(vm.STORE, s.Ident.Symbol)
	} else {
		out.bc.Instruction(vm.POP, nil)
	}
	return nil
}

func (out *compiler) compileBlockStmt(s *ast.BlockStmt, scope bool) error {
	if scope {
		out.bc.Instruction(vm.ENTER, nil)
	}

	for _, stmt := range s.Statements {
		if err := out.compileStmt(stmt); err != nil {
			return err
		}
	}

	if scope {
		out.bc.Instruction(vm.LEAVE, nil)
	}
	return nil
}

func (out *compiler) compileReturnStmt(s *ast.ReturnStmt) error {
	//// Push return expressions in reverse
	//for i := len(s.Returned) - 1; i >= 0; i-- {
	//	if err := out.compileExpr(s.Returned[i]); err != nil {
	//		return err
	//	}
	//}
	//
	//// Push int with amount of arguments
	//out.bc.Instruction(vm.PUSH, len(s.Returned))

	// TODO Multiple return values
	if len(s.Returned) > 0 {

		if err := out.compileExpr(s.Returned[0]); err != nil {
			return err
		}

	} else {
		out.bc.Instruction(vm.PUSH, nil)
	}
	out.bc.Instruction(vm.RET, nil)

	return nil
}

// compileForStmt compiles a for statement. This is by far the messiest implementation. TODO make it better.
func (out *compiler) compileForStmt(s *ast.ForStmt) error {
	out.bc.Instruction(vm.ENTER, nil)

	if s.Init != nil {
		if err := out.compileStmt(s.Init); err != nil {
			return err
		}
	}

	out.bc.Instruction(vm.ANCHOR, true)

	// BREAK
	skipBreakIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP, -1)
	endIndex := out.bc.Len()
	out.loopEnd.push(endIndex)
	out.bc.Instruction(vm.RESCUE, nil)
	endJumpIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP, -1)
	out.bc.SetArg(skipBreakIndex, out.bc.Len())

	// CONTINUE
	skipContinueIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP, -1)
	startIndex := out.bc.Len()
	out.loopBegin.push(startIndex)
	out.bc.Instruction(vm.RESCUE, nil)

	if s.Update != nil {
		if err := out.compileStmt(s.Update); err != nil {
			return err
		}
	}

	out.bc.SetArg(skipContinueIndex, out.bc.Len())

	var jumpIndex int
	if s.Cond != nil {
		if err := out.compileExpr(s.Cond); err != nil {
			return err
		}
		jumpIndex = out.bc.Len()
		out.bc.Instruction(vm.JUMP_F, nil)
	}

	//
	if err := out.compileStmt(s.Stmt); err != nil {
		return err
	}
	//

	out.bc.Instruction(vm.JUMP, startIndex)

	if s.Cond != nil {
		out.bc.SetArg(jumpIndex, out.bc.Len())
	}
	out.bc.SetArg(endJumpIndex, out.bc.Len())

	out.bc.Instruction(vm.ANCHOR, false)
	out.bc.Instruction(vm.LEAVE, nil)
	return nil
}

func (out *compiler) compileExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if err := out.compileBinaryExpr(e); err != nil {
			return err
		}
	case *ast.UnaryExpr:
		if err := out.compileUnaryExpr(e); err != nil {
			return err
		}
	case *ast.Number:
		if err := out.compileNumber(e); err != nil {
			return err
		}
	case *ast.Identifier:
		out.bc.Instruction(vm.LOAD, e.Symbol)
	case *ast.FunctionExpr:
		if err := out.compileFunctionExpr(e); err != nil {
			return err
		}
	case *ast.CallExpr:
		if err := out.compileCallExpr(e); err != nil {
			return err
		}
	case *ast.SubscriptExpr:
		if err := out.compileSubscriptExpr(e); err != nil {
			return err
		}
	case *ast.ArrayExpr:
		if err := out.compileArrayExpr(e); err != nil {
			return err
		}
	default:
		return ast.NewNodeError(expr, fmt.Sprintf("unknown expression type %T", expr))
	}
	return nil
}

func (out *compiler) compileArrayAssignStmt(s *ast.ArrayAssignStmt) error {
	if err := out.compileExpr(s.Expr); err != nil {
		return err
	}
	if err := out.compileExpr(s.Ident); err != nil {
		return err
	}
	if err := out.compileExpr(s.Index); err != nil {
		return err
	}
	out.bc.Instruction(vm.ARR_V, nil)
	return nil
}

func (out *compiler) compileNumber(e *ast.Number) error {
	f, err := strconv.ParseFloat(e.Value, 64)
	if err != nil {
		panic(err)
	}
	out.bc.Instruction(vm.PUSH, f)
	return nil
}

func (out *compiler) compileBinaryExpr(e *ast.BinaryExpr) error {
	// Check for operations that will short circuit
	switch e.Operator {
	case lexer.PIPE_PIPE:
		if err := out.compileExpr(e.Left); err != nil {
			return err
		}
		jumpTrueIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP_T, -1)

		// If not true
		if err := out.compileExpr(e.Right); err != nil {
			return err
		}
		secondTrueIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP_T, -1)

		// If false
		out.bc.Instruction(vm.PUSH, false) // <- will arrive here if false
		exitIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP, -1)

		// If true
		out.bc.SetArg(jumpTrueIndex, out.bc.Len())
		out.bc.SetArg(secondTrueIndex, out.bc.Len())
		out.bc.Instruction(vm.PUSH, true) // <- jump here if true
		out.bc.SetArg(exitIndex, out.bc.Len())

		return nil
	case lexer.AND_AND:
		if err := out.compileExpr(e.Left); err != nil {
			return err
		}
		jumpFalseIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP_F, -1)

		// If not true
		if err := out.compileExpr(e.Right); err != nil {
			return err
		}
		secondFalseIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP_F, -1)

		// If true
		out.bc.Instruction(vm.PUSH, true) // <- will arrive here if true
		exitIndex := out.bc.Len()
		out.bc.Instruction(vm.JUMP, -1)

		// If false
		out.bc.SetArg(jumpFalseIndex, out.bc.Len())
		out.bc.SetArg(secondFalseIndex, out.bc.Len())
		out.bc.Instruction(vm.PUSH, false) // <- jump here if false
		out.bc.SetArg(exitIndex, out.bc.Len())

		return nil
	default:
	}

	if err := out.compileExpr(e.Left); err != nil {
		return err
	}
	if err := out.compileExpr(e.Right); err != nil {
		return err
	}
	switch e.Operator {
	case lexer.PLUS:
		out.bc.Instruction(vm.ADD, nil)
	case lexer.MINUS:
		out.bc.Instruction(vm.SUB, nil)
	case lexer.ASTERISK:
		out.bc.Instruction(vm.MUL, nil)
	case lexer.SLASH:
		out.bc.Instruction(vm.DIV, nil)
	case lexer.EQUALS_EQUALS:
		out.bc.Instruction(vm.CMP, nil)
	case lexer.EXCLAMATION_EQUALS:
		out.bc.Instruction(vm.CMP, nil)
		out.bc.Instruction(vm.NOT, nil)
	case lexer.LESS_THAN:
		out.bc.Instruction(vm.CMP_LT, nil)
	case lexer.GREATER_THAN:
		out.bc.Instruction(vm.CMP_GT, nil)
	case lexer.LESS_THAN_EQUALS:
		out.bc.Instruction(vm.CMP_LTE, nil)
	case lexer.GREATER_THAN_EQUALS:
		out.bc.Instruction(vm.CMP_GTE, nil)
	default:
		return ast.NewNodeError(e, "unknown operator in binary expression")
	}

	return nil
}

func (out *compiler) compileUnaryExpr(e *ast.UnaryExpr) error {
	if err := out.compileExpr(e.Expr); err != nil {
		return err
	}

	switch e.Operator {
	case lexer.EXCLAMATION:
		out.bc.Instruction(vm.NOT, nil)
	case lexer.MINUS:
		out.bc.Instruction(vm.NEG, nil)
	case lexer.PLUS:
		// Do nothing
	default:
		return ast.NewNodeError(e, "unknown operator in unary expression")
	}

	return nil
}

func (out *compiler) compileFunctionExpr(e *ast.FunctionExpr) error {
	jumpIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP, -1)

	out.bc.Instruction(vm.ENTER, nil)

	argCountLabel := fmt.Sprintf("_argcount%d", out.bc.Len())

	if !e.IsVariadic {
		// Check if arg count matches exactly
		out.bc.Instruction(vm.PUSH, len(e.Params))
		out.bc.Instruction(vm.CMP, nil)
	} else {
		// Declare arg count
		out.bc.Instruction(vm.DECLARE, argCountLabel)

		out.bc.Instruction(vm.LOAD, argCountLabel)
		out.bc.Instruction(vm.PUSH, len(e.Params)-1) // Variadic args do accept an empty array.
		// Check if arg count matches or is greater
		out.bc.Instruction(vm.CMP_GTE, nil)
	}

	checkSuccessIndex := out.bc.Len()
	out.bc.Instruction(vm.JUMP_T, -1)

	// Return arg count error
	out.bc.Instruction(vm.PUSH, -1)
	out.compilePanic("arg count mismatch")

	out.bc.SetArg(checkSuccessIndex, out.bc.Len())

	for i, param := range e.Params {
		if e.IsVariadic && i >= len(e.Params)-1 {
			break
		}

		out.bc.Instruction(vm.DECLARE, param.Symbol)
	}

	if !e.IsVariadic {

	} else {
		index := len(e.Params) - 1

		// Calculate array size
		out.bc.Instruction(vm.LOAD, argCountLabel)
		out.bc.Instruction(vm.PUSH, len(e.Params)-1)
		out.bc.Instruction(vm.SUB, nil)

		out.bc.Instruction(vm.ARR_CR, nil)

		out.bc.Instruction(vm.DECLARE, e.Params[index].Symbol)
	}

	if err := out.compileBlockStmt(e.Body, false); err != nil {
		return err
	}

	out.bc.Instruction(vm.LEAVE, nil)

	// TODO Will need to inject the return statement later. It IS required.
	//out.bc.Instruction(vm.PUSH, 0) // Returning no values
	//
	//out.bc.Instruction(vm.PUSH, nil) // Return nil
	//out.bc.Instruction(vm.RET, "NOTHING")

	out.bc.SetArg(jumpIndex, out.bc.Len())

	// Push index of func (out *compiler)tion start. It is basically a pointer.

	function := vm.Func{
		Address: jumpIndex + 1,
	}

	out.bc.Instruction(vm.PUSH, function)

	return nil
}

func (out *compiler) compilePanic(s string) {
	out.bc.Instruction(vm.PANIC, s)
}

func (out *compiler) compileCallExpr(e *ast.CallExpr) error {
	// Push argument expressions in reverse to be declared in order in the call
	for i := len(e.Args) - 1; i >= 0; i-- {
		if err := out.compileExpr(e.Args[i]); err != nil {
			return err
		}
	}

	// Arg count
	out.bc.Instruction(vm.PUSH, len(e.Args))

	if err := out.compileExpr(e.Caller); err != nil {
		return err
	}

	frameReturnIndex := out.bc.Len()
	out.bc.Instruction(vm.FRAME, -1)

	out.bc.Instruction(vm.CALL, nil)
	out.bc.SetArg(frameReturnIndex, out.bc.Len())

	return nil
}

func (out *compiler) compileArrayExpr(e *ast.ArrayExpr) error {
	// Push argument expressions in reverse to be in order
	for i := len(e.Elements) - 1; i >= 0; i-- {
		if err := out.compileExpr(e.Elements[i]); err != nil {
			return err
		}
	}

	out.bc.Instruction(vm.PUSH, len(e.Elements))
	out.bc.Instruction(vm.ARR_CR, nil)
	return nil
}

func (out *compiler) compileSubscriptExpr(e *ast.SubscriptExpr) error {
	if err := out.compileExpr(e.Array); err != nil {
		return err
	}
	if err := out.compileExpr(e.Index); err != nil {
		return err
	}
	out.bc.Instruction(vm.ARR_ID, nil)
	return nil
}
