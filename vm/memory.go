package vm

func newScope(parent *Scope) *Scope {
	return &Scope{
		Parent:   parent,
		Declared: make(map[string]any),
	}
}

type Scope struct {
	Parent   *Scope
	Declared map[string]any
}

func (s *Scope) Set(name string, v any) {
	s.Declared[name] = v
}

func (s *Scope) Get(name string) any {
	if data, ok := s.Declared[name]; ok {
		return data
	}
	if s.Parent != nil {
		return s.Parent.Get(name)
	}
	return nil
}

type Data interface {
	Val() any
}

type Number struct {
	Value float64
}

func (n *Number) Val() any {
	return n.Value
}

type String struct {
	Value string
}

func (s *String) Val() any {
	return s.Value
}

func AsNumber(data Data) float64 {
	if n, ok := data.(*Number); ok {
		return n.Value
	}
	return 0
}

func AsString(data Data) string {
	if s, ok := data.(*String); ok {
		return s.Value
	}
	return ""
}
