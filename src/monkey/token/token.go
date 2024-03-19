package token

type TokenType int

const (
	Illegal TokenType = iota
	EOF
	Ident
	Int
	Assign
	Plus
	Minus
	Bang
	Asterisk
	Slash
	LessThan
	GreaterThan
	Equal
	NotEqual
	Comma
	Semicolon
	LParen
	RParen
	LSquirly
	RSquirly
	Function
	Let
	True
	False
	If
	Else
	Return
)

func (tt *TokenType) String() string {
	switch *tt {
	case Illegal:
		return "Illegal"
	case EOF:
		return "EOF"
	case Ident:
		return "Ident"
	case Int:
		return "Int"
	case Assign:
		return "Assign"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Bang:
		return "Bang"
	case Asterisk:
		return "Asterisk"
	case Slash:
		return "Slash"
	case LessThan:
		return "LessThan"
	case GreaterThan:
		return "GreaterThan"
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case LParen:
		return "LParen"
	case RParen:
		return "RParen"
	case LSquirly:
		return "LSquirly"
	case RSquirly:
		return "RSquirly"
	case Function:
		return "Function"
	case Let:
		return "Let"
	case True:
		return "True"
	case False:
		return "False"
	case If:
		return "If"
	case Else:
		return "Else"
	case Return:
		return "Return"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return Ident
}
