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
	_ = x[NEG-7]
	_ = x[CMP-8]
	_ = x[CMP_LT-9]
	_ = x[CMP_GT-10]
	_ = x[CMP_LTE-11]
	_ = x[CMP_GTE-12]
	_ = x[NOT-13]
	_ = x[DECLARE-14]
	_ = x[STORE-15]
	_ = x[LOAD-16]
	_ = x[JUMP-17]
	_ = x[JUMP_T-18]
	_ = x[JUMP_F-19]
	_ = x[JUMP_S-20]
	_ = x[ENTER-21]
	_ = x[LEAVE-22]
	_ = x[ECALL-23]
	_ = x[CALL-24]
	_ = x[RET-25]
	_ = x[ARR_CR-26]
	_ = x[ARR_ID-27]
	_ = x[ARR_V-28]
	_ = x[PANIC-29]
}

const _OpCode_name = "INVALIDPUSHPOPADDSUBMULDIVNEGCMPCMP_LTCMP_GTCMP_LTECMP_GTENOTDECLARESTORELOADJUMPJUMP_TJUMP_FJUMP_SENTERLEAVEECALLCALLRETARR_CRARR_IDARR_VPANIC"

var _OpCode_index = [...]uint8{0, 7, 11, 14, 17, 20, 23, 26, 29, 32, 38, 44, 51, 58, 61, 68, 73, 77, 81, 87, 93, 99, 104, 109, 114, 118, 121, 127, 133, 138, 143}

func (i OpCode) String() string {
	if i >= OpCode(len(_OpCode_index)-1) {
		return "OpCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}
