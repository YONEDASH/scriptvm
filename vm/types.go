package vm

type Func struct {
	Params []string
	Body   []Instr
}

func (f Func) value() {}
