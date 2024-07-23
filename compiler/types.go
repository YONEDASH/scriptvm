package compiler

import "script/vm"

var primitives = map[string]vm.TypeId{
	"int":   vm.Int,
	"float": vm.Float,
}
