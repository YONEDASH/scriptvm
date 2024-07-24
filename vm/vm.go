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

	vm.global.Declare("println", NewExternalFunc(func(v ...any) {
		fmt.Println(v...)
	}))

	vm.global.Declare("test", NewExternalFunc(func(do func() int) int {
		return do() * 2
	}))

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
		case PANIC:
			log.Fatalf("panic at %d: %v", i, instr.Arg)
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
	if v, err := Add(left, right); err != nil {
		vm.stack.Push(nil)
	} else {
		vm.stack.Push(v)
	}
}

func (vm *VM) sub() {
	left, right := vm.popBinary()
	if v, err := Sub(left, right); err != nil {
		vm.stack.Push(nil)
	} else {
		vm.stack.Push(v)
	}
}

func (vm *VM) mul() {
	left, right := vm.popBinary()
	if v, err := Mul(left, right); err != nil {
		vm.stack.Push(nil)
	} else {
		vm.stack.Push(v)
	}
}

func (vm *VM) div() {
	left, right := vm.popBinary()
	if v, err := Div(left, right); err != nil {
		vm.stack.Push(nil)
	} else {
		vm.stack.Push(v)
	}
}

func (vm *VM) neg() {
	if v, err := Neg(vm.stack.Pop()); err != nil {
		vm.stack.Push(nil)
	} else {
		vm.stack.Push(v)
	}
}

func (vm *VM) not() {
	vm.stack.Push(!vm.popBool())
}

func (vm *VM) cmp(code OpCode) {
	left, right := vm.popBinary()
	lType := TypeOf(left)
	rType := TypeOf(right)

	if lType != rType {
		log.Fatalf("cannot compare different types %v and %v", lType, rType)
	}

	switch lType {
	case Int:
		if leftInt, ok := left.(int); ok {
			if rightInt, ok := right.(int); ok {
				switch code {
				case CMP:
					vm.stack.Push(leftInt == rightInt)
				case CMP_LT:
					vm.stack.Push(leftInt < rightInt)
				case CMP_GT:
					vm.stack.Push(leftInt > rightInt)
				case CMP_LTE:
					vm.stack.Push(leftInt <= rightInt)
				case CMP_GTE:
					vm.stack.Push(leftInt >= rightInt)
				default:
					log.Fatalf("undefined comparison operation %v", code)
				}
				return
			}
			log.Fatalf("cannot compare integer with non-integer %v", code)
		}
	case Float:
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
					log.Fatalf("undefined comparison operation %v", code)
				}
				return
			}
			log.Fatalf("cannot compare number with non-number %v", code)
		}
	case Bool:
		if leftBool, ok := left.(bool); ok {
			if rightBool, ok := right.(bool); ok {
				switch code {
				case CMP:
					vm.stack.Push(leftBool == rightBool)
				default:
					log.Fatalf("undefined comparison operation %v", code)
				}
				return
			}
			log.Fatalf("cannot compare boolean with non-boolean %v", code)
		}
	default:
		log.Fatalf("undefined comparison for type %v of %v", lType, left)
	}
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
	case Type:
		vm.cast(t.Id)
		return
	case Func:
		address = t.Address
	case ExternalFunc:
		argCount := vm.stack.Pop().(int)
		result := t.Callback(vm, argCount)
		vm.stack.Push(result)
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

func (vm *VM) cast(t TypeId) {
	v := vm.stack.Pop()
	vt := TypeOf(v)

	switch t {
	case Any:
		vm.stack.Push(v)
		return
	case Array:
		if vt != Array {
			log.Fatalf("cannot cast to array from %v", vt)
		}
		vm.stack.Push(v.([]any))
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
