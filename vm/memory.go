package vm

func newFrame(parent *Frame) *Frame {
	return &Frame{
		Parent:   parent,
		Declared: make(map[string]any),
		start:    -1,
		end:      -1,
	}
}

type Frame struct {
	Parent   *Frame
	Declared map[string]any
	// end is the index of the instruction which invoked the function.
	start, end int
}

func (f *Frame) End() (*Frame, int) {
	if f.Parent != nil && f.end < 0 {
		return f.Parent.End()
	}

	return f, f.end
}

func (f *Frame) Assign(name string, v any) {
	if _, ok := f.Declared[name]; !ok {
		if f.Parent != nil {
			f.Parent.Assign(name, v)
		}
		return
	}
	f.Declared[name] = v
}

func (f *Frame) Declare(name string, v any) {
	f.Declared[name] = v
}

func (f *Frame) Get(name string) any {
	if data, ok := f.Declared[name]; ok {
		return data
	}
	if f.Parent != nil {
		return f.Parent.Get(name)
	}
	return nil
}
