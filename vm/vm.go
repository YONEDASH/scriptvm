package vm

import (
	"fmt"
	"log"
	"script"
	"strings"
)

func New() *VM {
	vm := &VM{
		stack:  make(Stack, 0, 16),
		global: newFrame(nil),
	}

	vm.global.Declare("int", Type{Int})
	vm.global.Declare("float", Type{Float})
	vm.global.Declare("bool", Type{Bool})

	return vm
}

type VM struct {
	global *Frame
	stack  Stack
}

func (vm *VM) Dump() string {
	dump := "## STACK ##\n"
	dump += script.Stringify(vm.stack) + "\n"

	dump += "## GLOBAL ##\n"
	dump += script.Stringify(vm.global) + "\n"

	return dump
}

const (
	debugInstructions = false
	debugStack        = false
)

func (vm *VM) Execute(bc Bytecode) error {
	if debugStack {
		fmt.Println(strings.TrimSpace(strings.ReplaceAll(script.Stringify(vm.stack), "\n", "")))
	}

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
			vm.stack.Push(instr.Arg)
		case POP:
			vm.stack.Pop()
		case ADD:
			vm.add()
		case SUB:
			vm.sub()
		case MUL:
			vm.mul()
		case DIV:
			vm.div()
		case CMP, CMP_LT, CMP_GT, CMP_LTE, CMP_GTE:
			vm.cmp(instr.Op)
		case NEG:
			vm.neg()
		case NOT:
			vm.not()
		case DECLARE:
			vm.declare(instr.Arg.(string))
		case LOAD:
			vm.load(instr.Arg.(string))
		case STORE:
			vm.store(instr.Arg.(string))
		case JUMP:
			i = instr.Arg.(int) - 1
		case JUMP_S:
			i = vm.stack.Pop().(int) - 1
		case JUMP_T:
			if vm.popBool() {
				i = instr.Arg.(int) - 1
			}
		case JUMP_F:
			if !vm.popBool() {
				i = instr.Arg.(int) - 1
			}
		case ENTER:
			vm.global = newFrame(vm.global)
		case LEAVE:
			if vm.global.Parent == nil {
				return fmt.Errorf("cannot leave global scope")
			}
			vm.global = vm.global.Parent
		case ECALL:
		//vm.call(instr.Arg)
		case CALL:
			vm.call(&i)
		case RET:
			var err error
			i, err = vm.ret(i)
			if err != nil {
				return err
			}
		case ARR_CR:
			vm.arrayCreate()
		case ARR_ID:
			vm.arrayIndex()
		case ARR_V:
			vm.arraySet()
		default:
			return fmt.Errorf("unknown opcode %vm in instruction %d", instr.Op, i)
		}
		if debugStack {
			fmt.Println(strings.TrimSpace(strings.ReplaceAll(script.Stringify(vm.stack), "\n", "")))
		}
	}
	return nil
}

func (vm *VM) ret(i int) (int, error) {
	if vm.global.Parent == nil {
		return 0, fmt.Errorf("cannot return without a frame")
	}
	p, index := vm.global.Origin()
	vm.global = p.Parent //return
	i = index - 1
	return i, nil
}

func (vm *VM) popBool() bool {
	value := vm.stack.Pop()

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

func (vm *VM) popBinary() (any, any) {
	right := vm.stack.Pop()
	left := vm.stack.Pop()
	return left, right
}

func (vm *VM) add() {
	left, right := vm.popBinary()
	vm.stack.Push(left.(float64) + right.(float64))
}

func (vm *VM) sub() {
	left, right := vm.popBinary()
	vm.stack.Push(left.(float64) - right.(float64))
}

func (vm *VM) mul() {
	left, right := vm.popBinary()
	vm.stack.Push(left.(float64) * right.(float64))
}

func (vm *VM) div() {
	left, right := vm.popBinary()
	vm.stack.Push(left.(float64) / right.(float64))
}

func (vm *VM) neg() {
	float := vm.stack.Pop().(float64)
	vm.stack.Push(-float)
}

func (vm *VM) not() {
	vm.stack.Push(!vm.popBool())
}

func (vm *VM) cmp(code OpCode) {
	left, right := vm.popBinary()

	if leftFloat, ok := left.(float64); ok {
		if rightFloat, ok := right.(float64); ok {
			switch code {
			case CMP:
				vm.stack.Push(leftFloat == rightFloat)
			case CMP_LT:
				vm.stack.Push(leftFloat < rightFloat)
			case CMP_GT:
				vm.stack.Push(leftFloat > rightFloat)
			case CMP_LTE:
				vm.stack.Push(leftFloat <= rightFloat)
			case CMP_GTE:
				vm.stack.Push(leftFloat >= rightFloat)
			default:
				log.Fatalf("undefined comparison operation %vm", code)
			}
			return
		}
		log.Fatalf("cannot compare number with non-number %vm", code)
	}

	if leftBool, ok := left.(bool); ok {
		if rightBool, ok := right.(bool); ok {
			switch code {
			case CMP:
				vm.stack.Push(leftBool == rightBool)
			default:
				log.Fatalf("undefined comparison operation %vm", code)
			}
			return
		}
		log.Fatalf("cannot compare boolean with non-boolean %vm", code)
	}

	log.Fatalf("cannot compare non-number with non-number %vm", code)
}

func (vm *VM) declare(s string) {
	vm.global.Declare(string(s), vm.stack.Pop())
}

func (vm *VM) load(key string) {
	val := vm.global.Get(string(key))
	vm.stack.Push(val)
}

func (vm *VM) store(s string) {
	vm.global.Assign(string(s), vm.stack.Pop())
}

func (vm *VM) arrayCreate() {
	size := vm.stack.Pop().(int)
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = vm.stack.Pop()
	}
	vm.stack.Push(arr)
}

func (vm *VM) arrayIndex() {
	index := int(vm.stack.Pop().(float64))
	arr := vm.stack.Pop().([]any)

	if index < 0 {
		index = len(arr) + index
	}

	if index < 0 || index >= len(arr) {
		vm.stack.Push(nil)
		return
	}

	vm.stack.Push(arr[index])
}

func (vm *VM) arraySet() {
	index := int(vm.stack.Pop().(float64))
	arr := vm.stack.Pop().([]any)

	if index < 0 {
		index = len(arr) + index
	}

	if index < 0 || index >= len(arr) {
		//vm.stack.Push(nil)
		return
	}

	arr[index] = vm.stack.Pop()
}

func (vm *VM) call(i *int) {
	top := vm.stack.Pop()

	address := -1

	switch t := top.(type) {
	case Func:
		address = t.Address
	case Type:
		vm.cast(t)
		return
	default:
		log.Fatalf("cannot call non-function %v", top)
		return
	}

	f := newFrame(vm.global)
	f.origin = *i + 1
	vm.global = f
	*i = address
}

func (vm *VM) cast(t Type) {
	v := vm.stack.Pop()
	vt := TypeOf(v)

	switch t.Id {
	case Int:
		switch vt {
		case Int:
			vm.stack.Push(v)
		case Float:
			vm.stack.Push(int(v.(float64)))
		case Bool:
			if v.(bool) {
				vm.stack.Push(int(1))
			} else {
				vm.stack.Push(int(0))
			}
		default:
			log.Fatalf("cannot cast to int from %v", vt)
		}
	case Float:
		switch vt {
		case Int:
			vm.stack.Push(float64(v.(int)))
		case Float:
			vm.stack.Push(v)
		case Bool:
			if v.(bool) {
				vm.stack.Push(float64(1))
			} else {
				vm.stack.Push(float64(0))
			}
		default:
			log.Fatalf("cannot cast to float from %v", vt)
		}
	case Bool:
		switch vt {
		case Int:
			vm.stack.Push(v.(int) != 0)
		case Float:
			vm.stack.Push(v.(float64) != 0)
		case Bool:
			vm.stack.Push(v)
		default:
			log.Fatalf("cannot cast to bool from %v", vt)
		}
	default:
		log.Fatalf("cannot cast to unknown type %vm", t)
	}
}
