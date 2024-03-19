package lexer

import (
	"testing"

	"monkey/token"

	"github.com/huandu/go-assert"
)

func TestNextToken(t *testing.T) {
	assert := assert.New(t)

	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LSquirly, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RSquirly, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.LessThan, "<"},
		{token.Int, "10"},
		{token.GreaterThan, ">"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LessThan, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LSquirly, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RSquirly, "}"},
		{token.Else, "else"},
		{token.LSquirly, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RSquirly, "}"},
		{token.Int, "10"},
		{token.Equal, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NotEqual, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	_, tokens := New(input)

	numTests := len(tests)

	index := 0
	for token := range tokens {
		assert.Equal(token.Type, tests[index].expectedType)

		assert.Equal(token.Literal, tests[index].expectedLiteral)

		index++
	}

	assert.Equal(index, numTests)
}
