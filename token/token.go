package token

import (
	"fmt"
	"log"
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

var keywordsReverse = map[Kind]string{
	Push:       "psh",
	Duplicate:  "dupl",
	Sum:        "sum",
	Divide:     "div",
	Subtract:   "sub",
	Multiply:   "mul",
	Jump:       "jmp",
	JumpIfTrue: "jif",
	Equal:      "eq",
}

func (t *Token) DetectMyKind() {
	value := t.IndentValue.Name

	if kind, ok := keywords[value]; ok {
		t.Kind = kind
		return
	}

	t.Kind = Identifier
}

func (t *Token) DetectMyString() string {
	if kind, ok := keywordsReverse[t.Kind]; ok {
		return kind
	}

	if t.Kind == Identifier {
		if t.IntegerValue != 0 {
			return fmt.Sprint(t.IntegerValue)
		}

		return t.StringValue
	}

	if t.Kind == Number {
		return fmt.Sprint(t.IntegerValue)
	}

	return ""
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

func (t *Token) DebugPrint() {
	log.Printf("	[Token]: [%s]\n", t.Kind)
	log.Printf("		IntegerValue: [%d]\n", t.IntegerValue)
	log.Printf("		StringValue: [%s]\n", t.StringValue)
	log.Printf("		LineStart: [%d]\n", t.LineStart)
	log.Printf("		LineEnd: [%d]\n", t.LineEnd)
	log.Printf("		CollumnStart: [%d]\n", t.CollumnStart)
	log.Printf("		CollumnEnd: [%d]\n", t.CollumnEnd)
	log.Printf("		[Indent]:\n")
	log.Printf("			Value: [%s]\n", t.Name)
	log.Printf("			Hash: [%d]\n", t.Hash)
}
