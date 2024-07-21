package vm

import "fmt"

type Stack []any

func (s *Stack) Push(val any) {
	fmt.Println("[STACK] PUSH", val)
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
	fmt.Println("[STACK] POP", val)
	return val
}
