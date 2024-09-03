package vm

import (
	"errors"
	"fmt"
)

//go:generate stringer -type=TypeId
type TypeId uint8

const (
	Invalid TypeId = iota
	Nil
	Any
	Int
	Float
	Bool
	Function
	Array
	ExternalFunction
)

func TypeOf(v any) TypeId {
	if v == nil {
		return Nil
	}
	switch v.(type) {
	case int:
		return Int
	case float64:
		return Float
	case bool:
		return Bool
	case Func:
		return Function
	case []any:
		return Array
	case any:
		return Any
	default:
		fmt.Printf("undefined type %T of %v\n", v, v)
		return Invalid
	}
}

type Type struct {
	Id TypeId
}

func (t Type) String() string {
	return fmt.Sprintf("<type %v>", t.Id)
}

type Func struct {
	//Params []string
	Address int
}

type ExternalFunc struct {
	Callback NativeFunc `json:"-"`
}

func (f Func) String() string {
	return fmt.Sprintf("<func %d>", f.Address)
}

var ErrTypeMismatch = errors.New("type mismatch")
var ErrTypeOperationUnsupported = errors.New("operation is unsupported for type")

func Add(a, b any) (any, error) {
	t := TypeOf(a)
	if t != TypeOf(b) {
		return nil, ErrTypeMismatch
	}
	switch t {
	case Int:
		return a.(int) + b.(int), nil
	case Float:
		return a.(float64) + b.(float64), nil
	default:
		return nil, ErrTypeOperationUnsupported
	}
}

func Sub(a, b any) (any, error) {
	t := TypeOf(a)
	if t != TypeOf(b) {
		return nil, ErrTypeMismatch
	}
	switch t {
	case Int:
		return a.(int) - b.(int), nil
	case Float:
		return a.(float64) - b.(float64), nil
	default:
		return nil, ErrTypeOperationUnsupported
	}
}

func Mul(a, b any) (any, error) {
	t := TypeOf(a)
	if t != TypeOf(b) {
		return nil, ErrTypeMismatch
	}
	switch t {
	case Int:
		return a.(int) * b.(int), nil
	case Float:
		return a.(float64) * b.(float64), nil
	default:
		return nil, ErrTypeOperationUnsupported
	}
}

func Div(a, b any) (any, error) {
	t := TypeOf(a)
	if t != TypeOf(b) {
		return nil, ErrTypeMismatch
	}
	switch t {
	case Int:
		return a.(int) / b.(int), nil
	case Float:
		return a.(float64) / b.(float64), nil
	default:
		return nil, ErrTypeOperationUnsupported
	}
}

func Neg(a any) (any, error) {
	t := TypeOf(a)
	switch t {
	case Int:
		return -a.(int), nil
	case Float:
		return -a.(float64), nil
	default:
		return nil, ErrTypeOperationUnsupported
	}
}
