package vm

func newFrame(parent *Frame) *Frame {
	return &Frame{
		Parent:   parent,
		Declared: make(map[string]any),
		origin:   -1,
	}
}

type Frame struct {
	Parent   *Frame
	Declared map[string]any
	// origin is the index of the instruction which invoked the function.
	origin int
}

func (f *Frame) Origin() (*Frame, int) {
	if f.Parent != nil && f.origin < 0 {
		return f.Parent.Origin()
	}

	return f, f.origin
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
