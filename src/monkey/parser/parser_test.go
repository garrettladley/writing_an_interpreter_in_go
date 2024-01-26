package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"strconv"
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

func TestIdentifierExpression(t *testing.T) {
	assert := assert.New(t)

	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	ident, ok := stmt.Expression.(*ast.Identifier)

	assert.Assert(ok)

	assert.Equal(ident.Value, "foobar")

	assert.Equal(ident.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	assert := assert.New(t)

	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	assert.Assert(ok)

	assert.Equal(literal.Value, int64(5))

	assert.Equal(literal.TokenLiteral(), "5")
}

func TestParsingPrefixExpressions(t *testing.T) {
	assert := assert.New(t)

	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		assert.Assert(ok)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		assert.Assert(ok)

		assert.Equal(exp.Operator, tt.operator)

		assert.Assert(testIntegerLiteral(assert, exp.Right, tt.value))
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	assert := assert.New(t)

	infixTests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		assert.Assert(ok)

		exp, ok := stmt.Expression.(*ast.InfixExpression)

		assert.Assert(ok)

		assert.Assert(testIntegerLiteral(assert, exp.Left, tt.leftVal))

		assert.Equal(exp.Operator, tt.operator)

		assert.Assert(testIntegerLiteral(assert, exp.Right, tt.rightVal))
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(program.String(), tt.expected)
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

func testIntegerLiteral(assert *assert.A, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	assert.Assert(ok)

	assert.Equal(integ.Value, value)
	assert.Equal(integ.TokenLiteral(), strconv.Itoa(int(value)))

	return true
}
