package lexer

import (
	"fmt"
	"script"
	"unicode"
)

type Buffer []rune

func (l *Buffer) Clear() {
	*l = (*l)[:0]
}

func (l *Buffer) Len() int {
	return len(*l)
}

func (l *Buffer) Append(r rune) {
	*l = append(*l, r)
}

func makeToken(pos int, id TokenId, lexeme string) Token {
	return Token{
		Pos:    pos,
		Id:     id,
		Lexeme: lexeme,
	}
}

type tokenizer struct {
	input  []rune
	pos    int
	buffer Buffer
	tokens []Token
	errors []error
}

// done Returns true if the tokenizer has reached EOF.
func (t *tokenizer) done() bool {
	return t.pos >= len(t.input)
}

// get Returns the rune at the given offset from the current position.
func (t *tokenizer) get(offset int) rune {
	p := t.pos + offset
	if p < 0 || p >= len(t.input) {
		return 0
	}
	return t.input[p]
}

// lex Returns the runes from the start to the end position. It is offset by the current position.
func (t *tokenizer) lex(start, end int) []rune {
	offset := t.pos
	start += offset
	end += offset
	if end < start || start < 0 || end > len(t.input) {
		return make([]rune, 0)
	}
	t.pos = end
	return t.input[start:end]
}

// push Pushes a new token to the tokens list at the current position. Will automatically push the buffer first using pushBuffer.
func (t *tokenizer) push(id TokenId, lexeme []rune) {
	t.pushBuffer()
	t.tokens = append(t.tokens, makeToken(t.pos, id, string(lexeme)))
}

// pushBuffer Pushes the current buffer to the tokens list at the current position minus buffer length.
func (t *tokenizer) pushBuffer() {
	l := t.buffer.Len()
	if l == 0 {
		return
	}
	start := t.pos - l
	lexeme := string(t.buffer)

	id := IDENTIFIER

	if keyword, ok := keywords[lexeme]; ok {
		id = keyword
	}

	t.buffer.Clear()
	t.tokens = append(t.tokens, Token{
		Pos:    start,
		Id:     id,
		Lexeme: lexeme,
	})
}

// number Pushes a number token to the tokens list.
func (t *tokenizer) number() {
	end := 1
	for unicode.IsDigit(t.get(end)) {
		end++
	}
	t.push(NUMBER, t.lex(0, end))
}

func Tokenize(input []byte) ([]Token, []error) {
	tr := &tokenizer{
		input:  []rune(string(input)),
		pos:    0,
		buffer: make(Buffer, 0, 64),
		tokens: make([]Token, 0, 1024),
		errors: make([]error, 0),
	}

	for !tr.done() {
		r := tr.get(0)
		switch r {
		case ' ', '\t', '\r':
			tr.pushBuffer()
			tr.pos++
		case '\n':
			tr.push(LF, tr.lex(0, 1))
		case '+':
			tr.push(PLUS, tr.lex(0, 1))
		case '-':
			tr.push(MINUS, tr.lex(0, 1))
		case '*':
			tr.push(ASTERISK, tr.lex(0, 1))
		case '/':
			tr.push(SLASH, tr.lex(0, 1))
		case '=':
			if tr.get(1) == '=' {
				tr.push(EQUALS_EQUALS, tr.lex(0, 2))
				continue
			}
			tr.push(EQUALS, tr.lex(0, 1))
		case ':':
			if tr.get(1) == '=' {
				tr.push(COLON_EQUALS, tr.lex(0, 2))
				continue
			}
			tr.push(COLON, tr.lex(0, 1))
		case '(':
			tr.push(OPEN_PAREN, tr.lex(0, 1))
		case ')':
			tr.push(CLOSE_PAREN, tr.lex(0, 1))
		case '{':
			tr.push(OPEN_BRACE, tr.lex(0, 1))
		case '}':
			tr.push(CLOSE_BRACE, tr.lex(0, 1))
		case '[':
			tr.push(OPEN_BRACKET, tr.lex(0, 1))
		case ']':
			tr.push(CLOSE_BRACKET, tr.lex(0, 1))
		case '<':
			if tr.get(1) == '=' {
				tr.push(LESS_THAN_EQUALS, tr.lex(0, 2))
				continue
			}
			tr.push(LESS_THAN, tr.lex(0, 1))
		case '>':
			if tr.get(1) == '=' {
				tr.push(GREATER_THAN_EQUALS, tr.lex(0, 2))
				continue
			}
			tr.push(GREATER_THAN, tr.lex(0, 1))
		case '!':
			if tr.get(1) == '=' {
				tr.push(EXCLAMATION_EQUALS, tr.lex(0, 2))
				continue
			}
			tr.push(EXCLAMATION, tr.lex(0, 1))
		case '^':
			tr.push(CIRCUMFLEX, tr.lex(0, 1))
		case '&':
			if tr.get(1) == '&' {
				tr.push(AND_AND, tr.lex(0, 2))
				continue
			}
			tr.push(AND, tr.lex(0, 1))
		case '|':
			if tr.get(1) == '|' {
				tr.push(PIPE_PIPE, tr.lex(0, 2))
				continue
			}
			tr.push(PIPE, tr.lex(0, 1))
		case ',':
			tr.push(COMMA, tr.lex(0, 1))
		case '.':
			if tr.get(1) == tr.get(2) && tr.get(2) == '.' {
				tr.push(DOT_DOT_DOT, tr.lex(0, 3))
			}
			tr.push(DOT, tr.lex(0, 1))
		default:
			if unicode.IsLetter(r) || (tr.buffer.Len() > 0 && unicode.IsDigit(r)) {
				tr.buffer.Append(r)
				tr.pos++
				continue
			}

			if unicode.IsDigit(r) {
				tr.number()
				continue
			}

			tr.errors = append(tr.errors, &script.PosError{
				Pos:     tr.pos,
				Message: fmt.Sprintf("unexpected character: %c", r),
			})
			tr.pos++
		}
	}
	tr.pushBuffer()
	tr.push(EOF, tr.lex(-1, -1))

	return tr.tokens, tr.errors
}
