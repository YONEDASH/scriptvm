// Code generated by "stringer -type=TypeId"; DO NOT EDIT.

package vm

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Invalid-0]
	_ = x[Nil-1]
	_ = x[Any-2]
	_ = x[Int-3]
	_ = x[Float-4]
	_ = x[Bool-5]
	_ = x[Function-6]
	_ = x[Array-7]
	_ = x[ExternalFunction-8]
}

const _TypeId_name = "InvalidNilAnyIntFloatBoolFunctionArrayExternalFunction"

var _TypeId_index = [...]uint8{0, 7, 10, 13, 16, 21, 25, 33, 38, 54}

func (i TypeId) String() string {
	if i >= TypeId(len(_TypeId_index)-1) {
		return "TypeId(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TypeId_name[_TypeId_index[i]:_TypeId_index[i+1]]
}
