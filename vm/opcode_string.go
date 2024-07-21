// Code generated by "stringer -type=OpCode"; DO NOT EDIT.

package vm

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INVALID-0]
	_ = x[PUSH-1]
	_ = x[POP-2]
	_ = x[ADD-3]
	_ = x[SUB-4]
	_ = x[MUL-5]
	_ = x[DIV-6]
	_ = x[LOAD-7]
	_ = x[STORE-8]
	_ = x[CALL-9]
	_ = x[RET-10]
}

const _OpCode_name = "INVALIDPUSHPOPADDSUBMULDIVLOADSTORECALLRET"

var _OpCode_index = [...]uint8{0, 7, 11, 14, 17, 20, 23, 26, 30, 35, 39, 42}

func (i OpCode) String() string {
	if i >= OpCode(len(_OpCode_index)-1) {
		return "OpCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}
