package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"monkey/token"
)

const eof = -1

type Lexer struct {
	input        string
	position     int // current position in input (points to current char)
	readPosition int // current reading position in input (after current char)
	width        int // width of last read char
	tokens       chan token.Token
}

func New(input string) (*Lexer, chan token.Token) {
	l := &Lexer{input: input, tokens: make(chan token.Token)}

	go l.run()

	return l, l.tokens
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func lex(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case isSpace(r):
			l.ignore()
		case r == '=':
			if l.peek() == '=' {
				l.next()
				l.emit(token.EQ)
			} else {
				l.emit(token.ASSIGN)
			}
		case r == '+':
			l.emit(token.PLUS)
		case r == '-':
			l.emit(token.MINUS)
		case r == '!':
			if l.peek() == '=' {
				l.next()
				l.emit(token.NOT_EQ)
			} else {
				l.emit(token.BANG)
			}
		case r == '/':
			l.emit(token.SLASH)
		case r == '*':
			l.emit(token.ASTERISK)
		case r == '<':
			l.emit(token.LT)
		case r == '>':
			l.emit(token.GT)
		case r == ';':
			l.emit(token.SEMICOLON)
		case r == ':':
			l.emit(token.COLON)
		case r == ',':
			l.emit(token.COMMA)
		case r == '(':
			l.emit(token.LPAREN)
		case r == ')':
			l.emit(token.RPAREN)
		case r == '{':
			l.emit(token.LSQUIRLY)
		case r == '}':
			l.emit(token.RSQUIRLY)
		case r == '[':
			l.emit(token.LBRACKET)
		case r == ']':
			l.emit(token.RBRACKET)
		case r == '"':
			return lexString
		case '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdent
		case r == eof:
			l.emit(token.EOF)
			return nil
		default:
			return l.errorf("unrecognized character in action: %#U", r)
		}
	}
}

func lexString(l *Lexer) stateFn {
	if l.next() == '"' {
		// handle the case where the string is empty
		l.readPosition = l.position
		l.emit(token.STRING)

		// move the position to the next character
		l.position += 1
		l.readPosition += 2

		return lex
	}

	for {
		switch r := l.next(); {
		case r == '"':
			// handle the " character at the beginning and end of the string
			startStringPosition := l.position + 1
			endStringPosition := l.readPosition - 1

			l.position = startStringPosition
			l.readPosition = endStringPosition

			l.emit(token.STRING)

			l.position = endStringPosition + 1
			l.readPosition = endStringPosition + 1

			return lex
		case r == eof:
			return l.errorf("unterminated string")
		}
	}
}

func lexNumber(l *Lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)

	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.position:l.readPosition])
	}

	l.emit(token.INT)

	return lex
}

func lexIdent(l *Lexer) stateFn {
	for isAlphaNumeric(l.next()) {
	}
	l.backup()

	l.emit(token.LookupIdent(l.input[l.position:l.readPosition]))

	return lex
}

func (l *Lexer) next() (r rune) {
	if l.readPosition >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.readPosition += l.width
	return r
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.readPosition -= l.width
}

func (l *Lexer) emit(t token.TokenType) {
	l.tokens <- token.Token{Type: t, Literal: l.input[l.position:l.readPosition]}
	l.position = l.readPosition
}

func (l *Lexer) ignore() {
	l.position = l.readPosition
}

func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{
		Type:    token.ILLEGAL,
		Literal: fmt.Sprintf(format, args...),
	}
	return nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
