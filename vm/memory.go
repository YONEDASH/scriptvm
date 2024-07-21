package vm

import "fmt"

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

func (s *Scope) Assign(name string, v any) {
	if _, ok := s.Declared[name]; !ok {
		if s.Parent != nil {
			s.Parent.Assign(name, v)
		}
		return
	}
	s.Declared[name] = v
	fmt.Println("[SCOPE] ASSIGN", name, v)
}

func (s *Scope) Declare(name string, v any) {
	s.Declared[name] = v
	fmt.Println("[SCOPE] DECLARE", name, v)
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
