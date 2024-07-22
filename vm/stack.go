package vm

type Stack []any

func (s *Stack) Push(val any) {
	*s = append(*s, val)
}

func (s *Stack) Top() any {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func (s *Stack) Pop() any {
	if len(*s) == 0 {
		return nil
	}
	val := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return val
}
