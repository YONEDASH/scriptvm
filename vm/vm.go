package vm

import (
	"fmt"
	"log"
	"script"
	"strings"
)

func New() *VM {
	vm := &VM{
		stack:  newStack(),
		cframe: newFrame(nil),
	}

	vm.cframe.Declare("int", Type{Int})
	vm.cframe.Declare("float", Type{Float})
	vm.cframe.Declare("bool", Type{Bool})

	vm.cframe.Declare("println", NewExternalFunc(func(v ...any) {
		fmt.Println(v...)
	}))

	return vm
}

type VM struct {
	cframe  *Frame
	stack   Stack
	pointer int
}

func (vm *VM) Err(msg string) {
	log.Fatalf("error at %d: %s", vm.pointer, msg)
}

func (vm *VM) Dump() string {
	//dump := "## STACK ##\n"
	//dump += script.Stringify(vm.stack) + "\n"

	dump := "## GLOBAL ##\n"
	dump += script.Stringify(vm.cframe) + "\n"

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

	vm.pointer = 0
	for ; vm.pointer < len(bc); vm.pointer++ {
		instr := bc[vm.pointer]
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
			vm.pointer = instr.Arg.(int) - 1
		case JUMP_T:
			if vm.popBool() {
				vm.pointer = instr.Arg.(int) - 1
			}
		case JUMP_F:
			if !vm.popBool() {
				vm.pointer = instr.Arg.(int) - 1
			}
		case ENTER:
			vm.cframe = newFrame(vm.cframe)
		case LEAVE:
			if vm.cframe.Parent == nil {
				return fmt.Errorf("cannot leave cframe scope")
			}
			vm.cframe = vm.cframe.Parent
		case CALL:
			vm.call(&vm.pointer)
		case RET:
			var err error
			vm.pointer, err = vm.ret(vm.pointer)
			if err != nil {
				return err
			}
		case ARR_INIT:
			vm.arrayInit()
		case ARR_CR:
			vm.arrayCreate()
		case ARR_ID:
			vm.arrayIndex()
		case ARR_V:
			vm.arraySet()
		case FRAME:
			vm.frame(vm.pointer, instr.Arg.(int))
		case ANCHOR:
			vm.cframe.anchor = instr.Arg.(bool)
		case RESCUE:
			vm.cframe = vm.cframe.Anchor()
		case JUMP_B:
			var err error
			vm.pointer, err = vm.jump_b(vm.pointer)
			if err != nil {
				return err
			}
		case PANIC:
			return fmt.Errorf("panic at %d: %v", vm.pointer, instr.Arg)
		default:
			return fmt.Errorf("unknown opcode %v in instruction %d", instr.Op, vm.pointer)
		}
		if debugStack {
			c := max(0, vm.stack.Cursor+1)
			fmt.Println(strings.TrimSpace(strings.ReplaceAll(script.Stringify(vm.stack.Array[:c]), "\n", "")))
		}
	}

	if vm.stack.Len() > 0 {
		return fmt.Errorf("memory leak: stack size is %d", vm.stack.Len())
	}

	return nil
}

func (vm *VM) ret(i int) (int, error) {
	if vm.cframe.Parent == nil {
		return 0, fmt.Errorf("cannot return without a frame")
	}
	p, index := vm.cframe.End()
	vm.cframe = p.Parent //return
	i = index - 1
	return i, nil
}

func (vm *VM) jump_b(i int) (int, error) {
	if vm.cframe.Parent == nil {
		return 0, fmt.Errorf("cannot jump_b without a frame")
	}
	p, _ := vm.cframe.End()
	vm.cframe = p //return to original frame
	i = p.start
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
		vm.Err(fmt.Sprintf("cannot compare different types %v and %v", lType, rType))
	}

	switch lType {
	case Nil:
		switch code {
		case CMP:
			vm.stack.Push(true)
		default:
			vm.Err(fmt.Sprintf("undefined comparison operation %v", code))
		}
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
					vm.Err(fmt.Sprintf("undefined comparison operation %v", code))
				}
				return
			}
			vm.Err(fmt.Sprintf("cannot compare integer with non-integer %v", code))
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
					vm.Err(fmt.Sprintf("undefined comparison operation %v", code))
				}
				return
			}
			vm.Err(fmt.Sprintf("cannot compare number with non-number %v", code))
		}
	case Bool:
		if leftBool, ok := left.(bool); ok {
			if rightBool, ok := right.(bool); ok {
				switch code {
				case CMP:
					vm.stack.Push(leftBool == rightBool)
				default:
					vm.Err(fmt.Sprintf("undefined comparison operation %v", code))
				}
				return
			}
			vm.Err(fmt.Sprintf("cannot compare boolean with non-boolean %v", code))
		}
	default:
		vm.Err(fmt.Sprintf("undefined comparison for type %v of %v", lType, left))
	}
}

func (vm *VM) declare(s string) {
	vm.cframe.Declare(string(s), vm.stack.Pop())
}

func (vm *VM) load(key string) {
	val := vm.cframe.Get(string(key))
	vm.stack.Push(val)
}

func (vm *VM) store(s string) {
	vm.cframe.Assign(string(s), vm.stack.Pop())
}

func (vm *VM) arrayInit() {
	size := vm.stack.Pop().(int)
	arr := make([]any, size)
	vm.stack.Push(arr)
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
	index := vm.stack.Pop().(int)
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
	index := vm.stack.Pop().(int)
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
		// Pop arg count
		vm.stack.Pop()
		// Cast
		vm.cast(t.Id)
		// Return
		*i, _ = vm.ret(*i)
		return
	case Func:
		address = t.Address
	case ExternalFunc:
		argCount := vm.stack.Pop().(int)
		result := t.Callback(vm, argCount)
		// Return
		*i, _ = vm.ret(*i)
		vm.stack.Push(result)
		return
	default:
		log.Fatalf("cannot call non-function %v @ %d", top, *i)
		return
	}

	*i = address
}

func (vm *VM) frame(current, end int) {
	f := newFrame(vm.cframe)
	f.start = current
	f.end = end
	vm.cframe = f
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
