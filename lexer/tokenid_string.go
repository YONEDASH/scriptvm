// Code generated by "stringer -type=TokenId"; DO NOT EDIT.

package lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INVALID-0]
	_ = x[EOF-1]
	_ = x[LF-2]
	_ = x[IDENTIFIER-3]
	_ = x[NUMBER-4]
	_ = x[STRING-5]
	_ = x[CHAR-6]
	_ = x[PLUS-7]
	_ = x[MINUS-8]
	_ = x[ASTERISK-9]
	_ = x[SLASH-10]
	_ = x[EQUALS-11]
	_ = x[COLON-12]
	_ = x[OPEN_PAREN-13]
	_ = x[CLOSE_PAREN-14]
	_ = x[OPEN_BRACE-15]
	_ = x[CLOSE_BRACE-16]
	_ = x[COLON_EQUALS-17]
	_ = x[EXCLAMATION-18]
	_ = x[EQUALS_EQUALS-19]
	_ = x[EXCLAMATION_EQUALS-20]
	_ = x[CIRCUMFLEX-21]
	_ = x[PIPE-22]
	_ = x[PIPE_PIPE-23]
	_ = x[AND-24]
	_ = x[AND_AND-25]
	_ = x[LESS_THAN-26]
	_ = x[GREATER_THAN-27]
	_ = x[LESS_THAN_EQUALS-28]
	_ = x[GREATER_THAN_EQUALS-29]
	_ = x[IF-30]
	_ = x[ELSE-31]
}

const _TokenId_name = "INVALIDEOFLFIDENTIFIERNUMBERSTRINGCHARPLUSMINUSASTERISKSLASHEQUALSCOLONOPEN_PARENCLOSE_PARENOPEN_BRACECLOSE_BRACECOLON_EQUALSEXCLAMATIONEQUALS_EQUALSEXCLAMATION_EQUALSCIRCUMFLEXPIPEPIPE_PIPEANDAND_ANDLESS_THANGREATER_THANLESS_THAN_EQUALSGREATER_THAN_EQUALSIFELSE"

var _TokenId_index = [...]uint16{0, 7, 10, 12, 22, 28, 34, 38, 42, 47, 55, 60, 66, 71, 81, 92, 102, 113, 125, 136, 149, 167, 177, 181, 190, 193, 200, 209, 221, 237, 256, 258, 262}

func (i TokenId) String() string {
	if i < 0 || i >= TokenId(len(_TokenId_index)-1) {
		return "TokenId(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenId_name[_TokenId_index[i]:_TokenId_index[i+1]]
}
