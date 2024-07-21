package vm

import (
	"fmt"
	"script"
)

func New() *VM {
	return &VM{
		stack:  make(Stack, 0, 16),
		global: newScope(nil),
	}
}

type VM struct {
	global *Scope
	stack  Stack
}

func (v *VM) Dump() string {
	dump := "## STACK ##\n"
	dump += script.Stringify(v.stack) + "\n"

	dump += "## GLOBAL ##\n"
	dump += script.Stringify(v.global) + "\n"

	return dump
}

func (v *VM) Execute(bc Bytecode) error {
	for i := 0; i < len(bc); i++ {
		instr := bc[i]
		switch instr.Op {
		case PUSH:
			v.stack.Push(instr.Arg)
		case POP:
			v.stack.Pop()
		case ADD:
			v.add()
		case SUB:
			v.sub()
		case MUL:
			v.mul()
		case DIV:
			v.div()
		case DECLARE:
			v.declare(instr.Arg.(string))
		case LOAD:
			v.load(instr.Arg.(string))
		case STORE:
			v.store(instr.Arg.(string))
		case ENTER:
			v.global = newScope(v.global)
		case LEAVE:
			if v.global.Parent == nil {
				return fmt.Errorf("cannot leave global scope")
			}
			v.global = v.global.Parent
		case CALL:
			v.call(instr.Arg)
		case RET:
			//return
		default:
			return fmt.Errorf("unknown opcode %v in instruction %d", instr.Op, i)
		}
	}
	return nil
}

func (v *VM) add() {
	left := v.stack.Pop()
	right := v.stack.Pop()
	v.stack.Push(left.(float64) + right.(float64))
}

func (v *VM) sub() {
	left := v.stack.Pop()
	right := v.stack.Pop()
	v.stack.Push(left.(float64) - right.(float64))
}

func (v *VM) mul() {
	left := v.stack.Pop()
	right := v.stack.Pop()
	v.stack.Push(left.(float64) * right.(float64))
}

func (v *VM) div() {
	left := v.stack.Pop()
	right := v.stack.Pop()
	v.stack.Push(left.(float64) / right.(float64))
}

func (v *VM) declare(s string) {
	v.global.Declare(s, v.stack.Pop())
}

func (v *VM) load(key string) {
	val := v.global.Get(key)
	v.stack.Push(val)
}

func (v *VM) store(s string) {
	v.global.Assign(s, v.stack.Pop())
}

func (v *VM) call(arg any) {
	if fn, ok := arg.(func(*VM)); ok {
		fn(v)
	} else {
		panic("not a function")
	}
}
