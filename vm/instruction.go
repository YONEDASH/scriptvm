package vm

import "fmt"

//go:generate stringer -type=OpCode
type OpCode uint8

const (
	INVALID OpCode = iota

	// PUSH <value>
	PUSH
	POP

	ADD
	SUB
	MUL
	DIV

	// NEG Negates the top of the stack.
	NEG

	CMP
	CMP_LT
	CMP_GT
	CMP_LTE
	CMP_GTE

	// NOT is equivalent to a logical NOT operation.
	NOT

	// DECLARE <name>
	DECLARE
	// STORE <name>
	STORE
	// LOAD <name>
	LOAD

	// JUMP <index>
	JUMP
	// JUMP_T <index> Will jump to given index if the top of the stack is true.
	JUMP_T
	// JUMP_F <index> Will jump to given index if the top of the stack is  false.
	JUMP_F
	// Deprecated: Will be removed.
	// JUMP_S Will jump to the index given on the top of the stack.
	JUMP_S

	// ENTER and LEAVE are used to create/exit a scope.
	ENTER
	LEAVE

	// CALL Used to call a function.
	CALL
	// FRAME <return_index> Used to initialize a new frame.
	FRAME
	// RET Return to ending position of frame and discard it.
	RET
	// JUMP_B Returns to the instruction after beginning of the last frame while keeping it. <=> RET
	JUMP_B

	// ARR_CR Creates an array. Array size on top of stack followed by elements.
	ARR_CR
	// ARR_ID Indexes into an array. TypeId will be pushed on top of stack.
	ARR_ID
	// ARR_V Sets an array element. TypeId on top of stack followed by index.
	ARR_V

	PANIC
)

type Bytecode []Instr

func (bc *Bytecode) Append(instr Instr) {
	*bc = append(*bc, instr)
}

func (bc *Bytecode) AppendBytecode(other Bytecode) {
	*bc = append(*bc, other...)
}

func (bc *Bytecode) Len() int {
	return len(*bc)
}

func (bc *Bytecode) Instruction(op OpCode, arg any) {
	bc.Append(Instr{Op: op, Arg: arg})
}

func (bc *Bytecode) SetArg(index int, arg any) {
	d := *bc
	d[index].Arg = arg
	*bc = d
}

func (bc *Bytecode) String() string {
	var s string
	for i, instr := range *bc {
		arg := ""
		if instr.Arg != nil {
			arg = fmt.Sprintf("%v", instr.Arg)
		}
		s += fmt.Sprintf("%3d\t%s\t%s\n", i, instr.Op, arg)
	}
	return s
}

type Instr struct {
	Op  OpCode
	Arg any
}
