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

	// ENTER and LEAVE are used to create/exit a scope.
	ENTER
	LEAVE

	CALL
	RET
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
