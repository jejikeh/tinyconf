package token

import (
	"unicode"

	"github.com/jejikeh/ambient/common"
)

type IndentValue struct {
	Name string
	Hash uint32
}

type TokenValue struct {
	IndentValue
	IntegerValue int
	FloatValue   float64
	StringValue  string
}

type Token struct {
	TokenValue

	Kind Kind

	LineStart    int
	CollumnStart int

	LineEnd    int
	CollumnEnd int
}

//go:generate stringer -type=TokenType
type Kind string

const (
	Push      = "PUSH"
	Duplicate = "DUPLICATE"

	Sum      = "SUM"
	Divide   = "DIV"
	Subtract = "SUB"
	Multiply = "MUL"

	Label = "LABEL"

	Comment = "COMMENT"

	Jump       = "JUMP"
	JumpIfTrue = "JUMP_IF_TRUE"

	Equal = "EQUAL"

	EndOfLine = "END_OF_FILE"

	Identifier = "IDENTIFIER"
	Number     = "NUMBER"
)

var keywords = map[string]Kind{
	"psh":  Push,
	"dupl": Duplicate,
	"sum":  Sum,
	"div":  Divide,
	"sub":  Subtract,
	"mul":  Multiply,
	"jmp":  Jump,
	"jif":  JumpIfTrue,
	"eq":   Equal,
}

func (t *Token) DetectMyKind() {
	value := t.IndentValue.Name

	if kind, ok := keywords[value]; ok {
		t.Kind = kind
		return
	}

	t.Kind = Identifier
}

func IsStartOfIdentifier(c rune) bool {
	if unicode.IsLetter(c) {
		return true
	}

	if c == '_' {
		return true
	}

	return false
}

func IsPartOfIdentifier(c rune) bool {
	if unicode.IsLetter(c) {
		return true
	}

	if unicode.IsDigit(c) {
		return true
	}

	if c == '_' {
		return true
	}

	return false
}

func IsPartOfNumber(c rune) bool {
	return unicode.IsDigit(c)
}

func IsPartOfFloat(c rune) bool {
	return unicode.IsDigit(c) || c == '.'
}

func IsWhitespace(c rune) bool {
	return unicode.IsSpace(c)
}

func (t *Token) SetIndentValue(s string) {
	common.AssertIfNot(len(s) > 0)

	t.IndentValue.Name = s
	t.IndentValue.Hash = common.Hash(s)
}
