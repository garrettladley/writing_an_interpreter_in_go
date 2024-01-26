package ast

import (
	"monkey/token"
	"testing"

	"github.com/huandu/go-assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(program.String(), "let myVar = anotherVar;")
}
