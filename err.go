package script

import (
	"fmt"
)

const PanicOnError = true

func NewPosError(pos int, message string) *PosError {
	e := &PosError{
		Pos:     pos,
		Message: message,
	}
	if PanicOnError {
		panic(e)
	}
	return e
}

type PosError struct {
	Pos     int
	Message string
}

func (e *PosError) Error() string {
	return fmt.Sprintf("get character %d: %s", e.Pos, e.Message)
}
