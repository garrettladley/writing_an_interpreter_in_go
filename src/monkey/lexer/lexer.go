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
				l.emit(token.Equal)
			} else {
				l.emit(token.Assign)
			}
		case r == '+':
			l.emit(token.Plus)
		case r == '-':
			l.emit(token.Minus)
		case r == '!':
			if l.peek() == '=' {
				l.next()
				l.emit(token.NotEqual)
			} else {
				l.emit(token.Bang)
			}
		case r == '/':
			l.emit(token.Slash)
		case r == '*':
			l.emit(token.Asterisk)
		case r == '<':
			l.emit(token.LessThan)
		case r == '>':
			l.emit(token.GreaterThan)
		case r == ';':
			l.emit(token.Semicolon)
		case r == ',':
			l.emit(token.Comma)
		case r == '(':
			l.emit(token.LParen)
		case r == ')':
			l.emit(token.RParen)
		case r == '{':
			l.emit(token.LSquirly)
		case r == '}':
			l.emit(token.RSquirly)
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

func lexNumber(l *Lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)

	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.position:l.readPosition])
	}

	l.emit(token.Int)

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
		Type:    token.Illegal,
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
