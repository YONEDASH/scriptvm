package vm

import "fmt"

type Bytecode []Instr

func (bc *Bytecode) Append(instr Instr) {
	*bc = append(*bc, instr)
}

func (bc *Bytecode) Instruction(op OpCode, arg any) {
	bc.Append(Instr{Op: op, Arg: arg})
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

//go:generate stringer -type=OpCode
type OpCode uint8

const (
	INVALID OpCode = iota

	PUSH
	POP

	ADD
	SUB
	MUL
	DIV

	LOAD
	STORE

	CALL
	RET
)
