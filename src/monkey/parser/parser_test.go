package parser

import (
	"strconv"
	"testing"

	"monkey/ast"
	"monkey/lexer"

	"github.com/huandu/go-assert"
)

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt := program.Statements[0]

		assert.Assert(testLetStatement(assert, stmt, tt.expectedIdentifier))

		val := stmt.(*ast.LetStatement).Value

		assert.Assert(testLiteralExpression(assert, val, tt.expectedValue))
	}
}

func TestReturnStatements(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt := program.Statements[0]

		returnStmt, ok := stmt.(*ast.ReturnStatement)

		assert.Assert(ok)

		assert.Equal(returnStmt.TokenLiteral(), "return")

		assert.Assert(testLiteralExpression(assert, returnStmt.ReturnValue, tt.expectedValue))
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
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		assert.Assert(testLiteralExpression(assert, exp.Right, tt.value))
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	assert := assert.New(t)

	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		assert.Assert(ok)

		assert.Assert(testInfixExpression(assert, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue))
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(program.String(), tt.expected)
	}
}

func TestBooleanExpression(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		assert.Equal(len(program.Statements), 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		assert.Assert(ok)

		boolean, ok := stmt.Expression.(*ast.Boolean)

		assert.Assert(ok)

		assert.Equal(boolean.Value, tt.expectedBoolean)
	}
}

func TestIfExpression(t *testing.T) {
	assert := assert.New(t)

	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)

	assert.Assert(ok)

	assert.Assert(testInfixExpression(assert, exp.Condition, "x", "<", "y"))

	assert.Equal(len(exp.Consequence.Statements), 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	assert.Assert(testIdentifier(assert, consequence.Expression, "x"))

	assert.Assert(exp.Alternative == nil)
}

func TestIfElseExpression(t *testing.T) {
	assert := assert.New(t)

	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)

	assert.Assert(ok)

	assert.Assert(testInfixExpression(assert, exp.Condition, "x", "<", "y"))

	assert.Equal(len(exp.Consequence.Statements), 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	assert.Assert(testIdentifier(assert, consequence.Expression, "x"))

	assert.Equal(len(exp.Alternative.Statements), 1)

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	assert.Assert(testIdentifier(assert, alternative.Expression, "y"))
}

func TestFunctionLiteralParsing(t *testing.T) {
	assert := assert.New(t)

	input := "fn(x, y) { x + y; }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	assert.Assert(ok)

	assert.Equal(len(function.Parameters), 2)

	assert.Assert(testLiteralExpression(assert, function.Parameters[0], "x"))

	assert.Assert(testLiteralExpression(assert, function.Parameters[1], "y"))

	assert.Equal(len(function.Body.Statements), 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	assert.Assert(testInfixExpression(assert, bodyStmt.Expression, "x", "+", "y"))
}

func TestFunctionParameterParsing(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		function := stmt.Expression.(*ast.FunctionLiteral)

		assert.Equal(len(function.Parameters), len(tt.expected))

		for i, ident := range tt.expected {
			testLiteralExpression(assert, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	assert := assert.New(t)

	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(assert, p)

	assert.Equal(len(program.Statements), 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	assert.Assert(ok)

	exp, ok := stmt.Expression.(*ast.CallExpression)

	assert.Assert(ok)

	assert.Assert(testIdentifier(assert, exp.Function, "add"))

	assert.Equal(len(exp.Arguments), 3)

	assert.Assert(testLiteralExpression(assert, exp.Arguments[0], 1))

	assert.Assert(testInfixExpression(assert, exp.Arguments[1], 2, "*", 3))

	assert.Assert(testInfixExpression(assert, exp.Arguments[2], 4, "+", 5))
}

func TestCallExpressionParameterParsing(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(assert, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp, ok := stmt.Expression.(*ast.CallExpression)

		assert.Assert(ok)

		assert.Assert(testIdentifier(assert, exp.Function, tt.expectedIdent))

		assert.Equal(len(exp.Arguments), len(tt.expectedArgs))

		for i, arg := range tt.expectedArgs {
			assert.Equal(exp.Arguments[i].String(), arg)
		}
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

func testIdentifier(assert *assert.A, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	assert.Assert(ok)

	assert.Equal(ident.Value, value)

	assert.Equal(ident.TokenLiteral(), value)

	return true
}

func testBooleanLiteral(assert *assert.A, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)

	assert.Assert(ok)

	assert.Equal(bo.Value, value)

	assert.Equal(bo.TokenLiteral(), strconv.FormatBool(value))

	return true
}

func testLiteralExpression(assert *assert.A, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(assert, exp, int64(v))
	case int64:
		return testIntegerLiteral(assert, exp, v)
	case string:
		return testIdentifier(assert, exp, v)
	case bool:
		return testBooleanLiteral(assert, exp, v)
	}

	assert.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testInfixExpression(assert *assert.A, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	assert.Assert(ok)

	assert.Assert(testLiteralExpression(assert, opExp.Left, left))

	assert.Equal(opExp.Operator, operator)

	assert.Assert(testLiteralExpression(assert, opExp.Right, right))

	return true
}
