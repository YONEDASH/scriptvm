package vm

import (
	"fmt"
	"log"
	"script"
	"strings"
)

func New() *VM {
	return &VM{
		stack:  make(Stack, 0, 16),
		global: newFrame(nil),
	}
}

type VM struct {
	global *Frame
	stack  Stack
}

func (v *VM) Dump() string {
	dump := "## STACK ##\n"
	dump += script.Stringify(v.stack) + "\n"

	dump += "## GLOBAL ##\n"
	dump += script.Stringify(v.global) + "\n"

	return dump
}

const (
	debugInstructions = false
	debugStack        = false
)

func (v *VM) Execute(bc Bytecode) error {
	for i := 0; i < len(bc); i++ {
		instr := bc[i]
		if debugInstructions {
			if instr.Arg != nil {
				fmt.Println(instr.Op, instr.Arg)
			} else {
				fmt.Println(instr.Op)
			}
		}
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
		case CMP, CMP_LT, CMP_GT, CMP_LTE, CMP_GTE:
			v.cmp(instr.Op)
		case NEG:
			v.neg()
		case NOT:
			v.not()
		case DECLARE:
			v.declare(instr.Arg.(string))
		case LOAD:
			v.load(instr.Arg.(string))
		case STORE:
			v.store(instr.Arg.(string))
		case JUMP:
			i = instr.Arg.(int) - 1
		case JUMP_S:
			i = v.stack.Pop().(int) - 1
		case JUMP_T:
			if v.popBool() {
				i = instr.Arg.(int) - 1
			}
		case JUMP_F:
			if !v.popBool() {
				i = instr.Arg.(int) - 1
			}
		case ENTER:
			v.global = newFrame(v.global)
		case LEAVE:
			if v.global.Parent == nil {
				return fmt.Errorf("cannot leave global scope")
			}
			v.global = v.global.Parent
		case CALL:
		//v.call(instr.Arg)
		case ICALL:
			f := newFrame(v.global)
			f.origin = i + 1
			v.global = f
			i = v.stack.Pop().(int) - 1
		case RET:
			if v.global.Parent == nil {
				return fmt.Errorf("cannot return without a frame")
			}
			p, index := v.global.Origin()
			v.global = p.Parent //return
			i = index - 1
		case ARR_CR:
			v.arrayCreate()
		case ARR_ID:
			v.arrayIndex()
		default:
			return fmt.Errorf("unknown opcode %v in instruction %d", instr.Op, i)
		}
		if debugStack {
			fmt.Println(strings.TrimSpace(strings.ReplaceAll(script.Stringify(v.stack), "\n", "")))
		}
	}
	return nil
}

func (v *VM) popBool() bool {
	value := v.stack.Pop()

	if value == nil {
		return false
	}

	switch t := value.(type) {
	case bool:
		return t
	case float64:
		// TODO remove this once booleans are implemented later
		return t != 0
	default:
		return true
	}
}

func (v *VM) popBinary() (any, any) {
	right := v.stack.Pop()
	left := v.stack.Pop()
	return left, right
}

func (v *VM) add() {
	left, right := v.popBinary()
	v.stack.Push(left.(float64) + right.(float64))
}

func (v *VM) sub() {
	left, right := v.popBinary()
	v.stack.Push(left.(float64) - right.(float64))
}

func (v *VM) mul() {
	left, right := v.popBinary()
	v.stack.Push(left.(float64) * right.(float64))
}

func (v *VM) div() {
	left, right := v.popBinary()
	v.stack.Push(left.(float64) / right.(float64))
}

func (v *VM) neg() {
	float := v.stack.Pop().(float64)
	v.stack.Push(-float)
}

func (v *VM) not() {
	v.stack.Push(!v.popBool())
}

func (v *VM) cmp(code OpCode) {
	left, right := v.popBinary()

	if leftFloat, ok := left.(float64); ok {
		if rightFloat, ok := right.(float64); ok {
			switch code {
			case CMP:
				v.stack.Push(leftFloat == rightFloat)
			case CMP_LT:
				v.stack.Push(leftFloat < rightFloat)
			case CMP_GT:
				v.stack.Push(leftFloat > rightFloat)
			case CMP_LTE:
				v.stack.Push(leftFloat <= rightFloat)
			case CMP_GTE:
				v.stack.Push(leftFloat >= rightFloat)
			default:
				log.Fatalf("undefined comparison operation %v", code)
			}
			return
		}
		log.Fatalf("cannot compare number with non-number %v", code)
	}

	if leftBool, ok := left.(bool); ok {
		if rightBool, ok := right.(bool); ok {
			switch code {
			case CMP:
				v.stack.Push(leftBool == rightBool)
			default:
				log.Fatalf("undefined comparison operation %v", code)
			}
			return
		}
		log.Fatalf("cannot compare boolean with non-boolean %v", code)
	}

	log.Fatalf("cannot compare non-number with non-number %v", code)
}

func (v *VM) declare(s string) {
	v.global.Declare(string(s), v.stack.Pop())
}

func (v *VM) load(key string) {
	val := v.global.Get(string(key))
	v.stack.Push(val)
}

func (v *VM) store(s string) {
	v.global.Assign(string(s), v.stack.Pop())
}

func (v *VM) arrayCreate() {
	size := v.stack.Pop().(int)
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = v.stack.Pop()
	}
	v.stack.Push(arr)
}

func (v *VM) arrayIndex() {
	index := int(v.stack.Pop().(float64))
	arr := v.stack.Pop().([]any)

	if index < 0 || index > len(arr) {
		v.stack.Push(nil)
		return
	}

	v.stack.Push(arr[index])
}
