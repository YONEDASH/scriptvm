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
	_ = x[COMMA-13]
	_ = x[OPEN_PAREN-14]
	_ = x[CLOSE_PAREN-15]
	_ = x[OPEN_BRACE-16]
	_ = x[CLOSE_BRACE-17]
	_ = x[OPEN_BRACKET-18]
	_ = x[CLOSE_BRACKET-19]
	_ = x[COLON_EQUALS-20]
	_ = x[EXCLAMATION-21]
	_ = x[EQUALS_EQUALS-22]
	_ = x[EXCLAMATION_EQUALS-23]
	_ = x[CIRCUMFLEX-24]
	_ = x[PIPE-25]
	_ = x[PIPE_PIPE-26]
	_ = x[AND-27]
	_ = x[AND_AND-28]
	_ = x[LESS_THAN-29]
	_ = x[GREATER_THAN-30]
	_ = x[LESS_THAN_EQUALS-31]
	_ = x[GREATER_THAN_EQUALS-32]
	_ = x[IF-33]
	_ = x[ELSE-34]
	_ = x[RETURN-35]
}

const _TokenId_name = "INVALIDEOFLFIDENTIFIERNUMBERSTRINGCHARPLUSMINUSASTERISKSLASHEQUALSCOLONCOMMAOPEN_PARENCLOSE_PARENOPEN_BRACECLOSE_BRACEOPEN_BRACKETCLOSE_BRACKETCOLON_EQUALSEXCLAMATIONEQUALS_EQUALSEXCLAMATION_EQUALSCIRCUMFLEXPIPEPIPE_PIPEANDAND_ANDLESS_THANGREATER_THANLESS_THAN_EQUALSGREATER_THAN_EQUALSIFELSERETURN"

var _TokenId_index = [...]uint16{0, 7, 10, 12, 22, 28, 34, 38, 42, 47, 55, 60, 66, 71, 76, 86, 97, 107, 118, 130, 143, 155, 166, 179, 197, 207, 211, 220, 223, 230, 239, 251, 267, 286, 288, 292, 298}

func (i TokenId) String() string {
	if i < 0 || i >= TokenId(len(_TokenId_index)-1) {
		return "TokenId(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenId_name[_TokenId_index[i]:_TokenId_index[i+1]]
}
