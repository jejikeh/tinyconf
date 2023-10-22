package lexer

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/token"
)

// @IMPORTANT: The InputCursor is the current position in the InputSource - 1

type Lexer struct {
	CurrentLineNumber                  int
	CurrentLineExternalFileErrorReport int
	CurrentLineCharacterIndex          int

	Tokens      []token.Token
	TokenCursor int

	InputSource []rune
	InputCursor int

	TotalLinesProcessed int
}

func NewLexer(source string) *Lexer {
	return &Lexer{}
}

// eatCharacter just increments the InputCursor by 1
//
// Also increments the CurrentLineCharacterIndex and another
// line counters for error reporting and other things
func (l *Lexer) eatCharacter() {
	if l.InputCursor >= len(l.InputSource) {
		return
	}

	if l.InputSource[l.InputCursor] == '\n' {
		l.CurrentLineNumber++
		l.TotalLinesProcessed++
		l.CurrentLineCharacterIndex = 0
	}

	l.InputCursor++
	l.CurrentLineCharacterIndex++
}

func (l *Lexer) peekNextCharacter() (rune, error) {
	if l.InputCursor >= len(l.InputSource) {
		return -1, fmt.Errorf("unexpected end of file")
	}

	return l.InputSource[l.InputCursor], nil
}

func (l *Lexer) getDecimalDigit() (int, error) {
	c, err := l.peekNextCharacter()
	if err != nil {
		return 0, err
	}

	if unicode.IsDigit(c) {
		return int(c - '0'), nil
	}

	return 0, fmt.Errorf("expected decimal digit, but got unexpected character: [%v]", c)
}

func (l *Lexer) setStartOfToken(t *token.Token) {
	t.LineStart = l.CurrentLineNumber
	t.CollumnStart = l.CurrentLineCharacterIndex
}

func (l *Lexer) setEndOfToken(t *token.Token) {
	t.LineEnd = l.CurrentLineNumber
	t.CollumnEnd = l.CurrentLineCharacterIndex
}

func (l *Lexer) throwBackOneCharacter() {
	common.AssertIfNot(l.InputCursor > 0)

	l.InputCursor--
	l.CurrentLineCharacterIndex--
}

func (l *Lexer) parseIndentifier(t *token.Token) error {
	strBuilder := strings.Builder{}

	c, err := l.peekNextCharacter()
	if err != nil {
		return err
	}

	if token.IsPartOfIdentifier(c) {
		for {
			c, err = l.peekNextCharacter()
			if err != nil {
				return err
			}

			if token.IsPartOfIdentifier(c) {
				l.eatCharacter()
				strBuilder.WriteRune(c)
				continue
			}

			break
		}
	} else {
		return fmt.Errorf("expected part of identifier, but got unexpected character: [%v]", c)
	}

	t.SetIndentValue(strBuilder.String())
	l.setEndOfToken(t)

	return nil
}

func (l *Lexer) makeIdentifierOrKeyword() (token.Token, error) {
	// @todo: get_unused_keyword()
	t := &token.Token{}
	t.Kind = token.Indentifier

	err := l.parseIndentifier(t)
	if err != nil {
		return *t, err
	}

	t.DetectMyKind()

	return *t, nil
}

func (l *Lexer) makeNumber() (token.Token, error) {
	// @todo: get_unused_keyword()
	t := &token.Token{}
	t.Kind = token.Number

	strBuilder := strings.Builder{}

	c, err := l.peekNextCharacter()
	if err != nil {
		return *t, err
	}

	if token.IsPartOfNumber(c) {
		// @todo: getDecimalDigit()
		for {
			c, err = l.peekNextCharacter()
			if err != nil {
				return *t, err
			}

			if token.IsPartOfNumber(c) {
				l.eatCharacter()
				strBuilder.WriteRune(c)
				continue
			}

			break
		}
	} else {
		return *t, fmt.Errorf("expected start of number, but got unexpected character: [%v]", c)
	}

	num, err := strconv.Atoi(strBuilder.String())
	if err != nil {
		return *t, err
	}

	t.IntegerValue = num

	return *t, nil
}

func (l *Lexer) eatInputDueToBlockComment() error {
	commentDepth := 1

	for commentDepth > 0 {
		c, err := l.peekNextCharacter()
		if err != nil {
			return err
		}

		if c == '/' {
			l.eatCharacter()
			c, err = l.peekNextCharacter()
			if err != nil {
				return err
			}

			if c == '*' {
				l.eatCharacter()
				commentDepth++
			}
		} else if c == '*' {
			l.eatCharacter()

			c, err = l.peekNextCharacter()
			if err != nil {
				return err
			}

			if c == '/' {
				l.eatCharacter()
				commentDepth--
			}
		} else {
			l.eatCharacter()
		}
	}

	return nil
}

func (l *Lexer) eatUntilNewLine() {
	for {
		c, err := l.peekNextCharacter()
		if err != nil {
			return
		}

		if c == '\n' {
			return
		}

		l.eatCharacter()
	}
}

func (l *Lexer) composeNewToken() (token.Token, error) {
	t := &token.Token{}

	for {
		c, err := l.peekNextCharacter()
		if err != nil {
			if c == -1 {
				t.Kind = token.EndOfLine
				return *t, nil
			}

			return *t, err
		}

		for token.IsWhitespace(c) {
			l.eatCharacter()
			c, err = l.peekNextCharacter()
			if err != nil {
				return *t, err
			}
		}

		if token.IsStartOfIdentifier(c) {
			l.setStartOfToken(t)
			return l.makeIdentifierOrKeyword()
		}

		if token.IsPartOfNumber(c) {
			return l.makeNumber()
		} else {
			return *t, fmt.Errorf("unexpected character: [%v] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
		}
	}
}

func NewLexerFromFile(filepath string) *Lexer {
	content, err := os.ReadFile(filepath)
	if err != nil {
		color.Set(color.FgRed)
		defer color.Unset()
		log.Fatalf("Error reading file: [%v]\n", err)
	}

	return &Lexer{
		InputSource: []rune(string(content)),
	}
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := []token.Token{}

	for l.InputCursor < len(l.InputSource) {
		t, err := l.composeNewToken()
		if err != nil {
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Printf("Error: %s\n", err)
		}

		log.Printf("Token: [%s]\n", t.Kind)

		tokens = append(tokens, t)
	}

	return tokens
}

func PrintDebugTokens(tokens []token.Token) {
	log.Println("Tokens:")
	for i, token := range tokens {
		log.Printf("[%d] ", i)
		logDebugToken(&token)
	}
}

func logDebugToken(t *token.Token) {
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

// v1 lexer :(
/*


func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}
	token := &Token{}

	for l.Cursor < len(l.Source) {
		l.lexWhitespace()

		token = l.lexSyntaxToken()
		if token != nil {
			tokens = append(tokens, *token)
			continue
		}

		token = l.lexOperandTokens()
		if token != nil {
			tokens = append(tokens, *token)
			continue
		}

		token = l.lexNumberTokens()
		if token != nil {
			tokens = append(tokens, *token)
			continue
		}

		color.Set(color.FgHiRed)
		defer color.Unset()

		log.Printf("Unknown token:  [%s:%d]", string(l.Source[l.Cursor]), l.Cursor)
		l.Cursor++
	}

	return tokens
}

func (l *Lexer) lexWhitespace() *Token {
	for l.Cursor < len(l.Source) {
		symbol := l.Source[l.Cursor]
		if unicode.IsSpace(symbol) {
			l.Cursor++
			continue
		}

		break
	}

	return nil
}

func (l *Lexer) lexSyntaxToken() *Token {
	input := ""
	save_cursor := l.Cursor

	for l.Cursor < len(l.Source) && !unicode.IsSpace(l.Source[l.Cursor]) {
		input += string(l.Source[l.Cursor])
		l.Cursor++
	}

	// NOTE(jejikeh): If we find a keyword, then maybe we can wait until space. In case we have a keywords
	// which have similar spelling with other keywords.
	if t, ok := Keywords[input]; ok {
		l.Cursor++
		return &Token{
			Value:    input,
			Kind:     t,
			Location: save_cursor,
		}
	}

	if t, ok := Keywords[string(input[0])]; ok {
		if t == Comment {
			return &Token{
				Value:    input,
				Kind:     t,
				Location: save_cursor,
			}
		}
	}

	l.Cursor = save_cursor
	return nil
}

func (l *Lexer) lexOperandTokens() *Token {
	token := []rune{}
	save_cursor := l.Cursor

	for l.Cursor < len(l.Source) && !unicode.IsSpace(l.Source[l.Cursor]) {
		symbol := l.Source[l.Cursor]
		if unicode.IsLetter(symbol) || (len(token) != 0 && unicode.IsDigit(symbol)) {
			token = append(token, symbol)
			l.Cursor++
			continue
		}

		break
	}

	if len(token) == 0 {
		l.Cursor = save_cursor
		return nil
	}

	return &Token{
		Value:    string(token),
		Kind:     OperandString,
		Location: save_cursor,
	}
}

// NOTE(jejikeh): for now we have only operands. Maybe just for lexing its okay.
func (l *Lexer) lexNumberTokens() *Token {
	token := []rune{}
	save_cursor := l.Cursor

	for l.Cursor < len(l.Source) && !unicode.IsSpace(l.Source[l.Cursor]) {
		symbol := l.Source[l.Cursor]
		if unicode.IsDigit(symbol) {
			token = append(token, symbol)
			l.Cursor++
			continue
		}

		break
	}

	if len(token) == 0 {
		l.Cursor = save_cursor
		return nil
	}

	return &Token{
		Value:    string(token),
		Kind:     OperandNumber,
		Location: l.Cursor - 1,
	}
}

*/
