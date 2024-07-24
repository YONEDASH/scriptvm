package vm

const MaxStackSize = 512

type Stack struct {
	Array  [MaxStackSize]any
	Cursor int
}

func (s *Stack) Push(val any) {
	s.Cursor++
	s.Array[s.Cursor] = val
}

func (s *Stack) Top() any {
	if s.Cursor < 0 {
		return nil
	}
	return s.Array[s.Cursor]
}

func (s *Stack) Pop() any {
	if s.Cursor < 0 {
		return nil
	}

	val := s.Array[s.Cursor]
	s.Cursor--
	return val
}

func (s *Stack) Len() int {
	return s.Cursor
}

func newStack() Stack {
	return Stack{
		Array:  [MaxStackSize]any{},
		Cursor: -1,
	}
}
