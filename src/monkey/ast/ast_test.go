package ast

import (
	"testing"

	"monkey/token"

	"github.com/huandu/go-assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.Ident, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Ident, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(program.String(), "let myVar = anotherVar;")
}
