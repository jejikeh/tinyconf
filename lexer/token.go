package lexer

//go:generate stringer -type=TokenType
type TokenType int

const (
	// Math
	Sum TokenType = iota
	Div
	Sub
	Mul

	// Memory
	Push
	Duplicate

	// Control
	Jump
	JumpIfTrue
	Equal

	// Misc
	LabelDeclaration

	OperandString
	OperandNumber

	Comment
)

type Token struct {
	Value    string
	Kind     TokenType
	Location int
}

var Keywords = map[string]TokenType{
	// Language Semantics
	"push":  Push,
	"dupl":  Duplicate,
	"sum":   Sum,
	"div":   Div,
	"sub":   Sub,
	"mul":   Mul,
	"jump":  Jump,
	"jtrue": JumpIfTrue,
	"equal": Equal,

	// Special Symbols
	":": LabelDeclaration,
	"#": Comment,
}
