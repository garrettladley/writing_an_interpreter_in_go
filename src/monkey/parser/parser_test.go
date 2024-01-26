package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"

	"github.com/huandu/go-assert"
)

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.NotEqual(program, nil)

	assert.Equal(len(program.Statements), 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		assert.Assert(testLetStatement(assert, program.Statements[i], tt.expectedIdentifier))
	}
}

func TestReturnStatements(t *testing.T) {
	assert := assert.New(t)
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.NotEqual(program, nil)

	assert.Equal(len(program.Statements), 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		assert.Assert(ok)

		assert.Equal(returnStmt.TokenLiteral(), "return")
	}
}

func checkParserErrors(assert *assert.A, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	assert.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		assert.Errorf("parser error: %q", msg)
	}

	assert.FailNow()
}

func testLetStatement(assert *assert.A, s ast.Statement, name string) bool {
	assert.Equal(s.TokenLiteral(), "let")

	letStmt, ok := s.(*ast.LetStatement)
	assert.Assert(ok)

	assert.Equal(letStmt.Name.Value, name)
	assert.Equal(letStmt.Name.TokenLiteral(), name)

	return true
}
