package vm

import (
	"fmt"
	"reflect"
)

type NativeFunc func(vm *VM, argCount int) any

func (n NativeFunc) String() string {
	return "<native function>"
}

func NewExternalFunc(f any) ExternalFunc {
	fn, err := newExternalCallback(f)
	if err != nil {
		panic(err)
	}
	return ExternalFunc{fn}
}

func newExternalCallback(f any) (NativeFunc, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("%v is not a function", f)
	}

	if t.NumOut() > 1 {
		return nil, fmt.Errorf("%v cannot return more than 1 value", f)
	}

	v := reflect.ValueOf(f)

	numIn := t.NumIn()
	types := make([]TypeId, numIn)

	variadicIndex := -1
	if t.IsVariadic() {
		variadicIndex = numIn - 1
	}

	for i := 0; i < numIn; i++ {
		in := t.In(i)

		var val any

		if variadicIndex == i {
			if in.Kind() != reflect.Slice {
				return nil, fmt.Errorf("variadic argument must be a slice")
			}

			val = reflect.New(in).Interface()
		} else {
			val = reflect.New(in).Elem().Interface()
		}

		vmType := TypeOf(val)

		fmt.Println(i, vmType, in)
		fmt.Printf("type=%T valtype=%T\n", in, val)

		if vmType != Invalid {
			types[i] = vmType
			continue
		}

		switch in.Kind() {
		case reflect.Func:
			types[i] = Function
		default:
			return nil, fmt.Errorf("unsupported argument type at position %d: %v", i, in.Kind())
		}
	}

	return func(vm *VM, argCount int) any {
		maxArgs := min(numIn, argCount)
		if variadicIndex >= 0 {
			maxArgs = argCount
		}

		args := make([]any, maxArgs)
		for i := 0; i < maxArgs; i++ {
			var argType TypeId
			if i < numIn {
				argType = types[i]
			} else {
				argType = types[variadicIndex]
			}
			vm.cast(argType)
			args[i] = vm.stack.Pop()
		}

		if argCount > maxArgs {
			for i := maxArgs; i < argCount; i++ {
				// Pop away unused arguments.
				vm.stack.Pop()
			}
		}

		values := make([]reflect.Value, len(args))

		for i := 0; i < argCount; i++ {
			values[i] = reflect.ValueOf(args[i])
		}

		v.Call(values)

		return nil
	}, nil
}
