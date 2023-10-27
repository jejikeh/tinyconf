package lexer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	return &Lexer{InputSource: []rune(source)}
}

func NewLexerFromSource(filepath string) *Lexer {
	content, err := os.ReadFile(filepath)
	if err != nil {
		color.Set(color.FgHiRed)
		defer color.Unset()
		log.Fatalf("Error reading file: [%v]\n", err)
	}

	return &Lexer{
		InputSource: []rune(string(content)),
	}
}

func NewLexerFromBinary(filepath string) *Lexer {
	return &Lexer{
		Tokens: loadFromBinary(filepath),
	}
}

func (l *Lexer) DumpTokensToBinary(outputPath string) {
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputPath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err = enc.Encode(l.Tokens)

	if err != nil {
		log.Fatal("Error encoding instructions: ", err)
	}

	_, err = f.Write(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func loadFromBinary(sourcePath string) []token.Token {
	var tokens []token.Token

	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(content))
	err = dec.Decode(&tokens)
	if err != nil {
		log.Fatal("Error decoding instructions: ", err)
	}

	return tokens
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := []token.Token{}

	for {
		t, err := l.composeNewToken()
		if err != nil {
			color.Set(color.FgHiRed)
			defer color.Unset()

			log.Fatalf("Error: %s\n", err)
		}

		tokens = append(tokens, t)

		if t.Kind == token.EndOfLine {
			break
		}
	}

	tokens = l.resolveLabelIdentifierDeclaration(tokens)

	return tokens
}

func PrintDebugTokens(tokens []token.Token) {
	log.Println("Tokens:")
	for i, token := range tokens {
		log.Printf("[%d] ", i)
		logDebugToken(&token)
	}
}

// @todo: refactor this please
func (l *Lexer) resolveLabelIdentifierDeclaration(tokens []token.Token) []token.Token {
	labelIndexs := make(map[string]int)
	newTokens := []token.Token{}

	// @todo: remove two loops
	for i, t := range tokens {
		if t.Kind == token.Label {
			labelIndexs[t.Name] = i
		}
	}

	for _, t := range tokens {
		if t.Kind == token.Identifier {
			if label, ok := labelIndexs[t.Name]; ok {
				t.IntegerValue = label
			} else {
				color.Set(color.FgHiRed)
				defer color.Unset()
				log.Printf("Label [%s] (%d:%d) not found! Default to -1\n", t.Name, t.LineStart, t.LineEnd)
				t.IntegerValue = -1
			}
		}

		newTokens = append(newTokens, t)
	}

	return newTokens
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
		return -1, fmt.Errorf(`unexpected end of file!
		Current cursor: [%d]
		Source length: [%d]
		Current line: [%d]
		Current column: [%d]`, l.InputCursor, len(l.InputSource), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
	}

	return l.InputSource[l.InputCursor], nil
}

func (l *Lexer) peekNextCharacterIgnoringRunes(r rune) (rune, error) {
	c, err := l.peekNextCharacter()
	if err != nil {
		return -1, err
	}

	if c == r {
		l.eatCharacter()
		return l.peekNextCharacterIgnoringRunes(r)
	}

	return c, nil
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

	if token.IsStartOfIdentifier(c) {
		for {
			c, err = l.peekNextCharacter()
			if err != nil {
				if strBuilder.Len() > 0 {
					break
				}

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
	t.Kind = token.Identifier

	l.setStartOfToken(t)

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

	l.setStartOfToken(t)

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
				if strBuilder.Len() > 0 {
					break
				}

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
		return *t, fmt.Errorf("expected start of number, but got unexpected character: [%s] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
	}

	t.SetIndentValue(strBuilder.String())
	num, err := strconv.Atoi(t.Name)
	if err != nil {
		return *t, err
	}

	t.IntegerValue = num
	l.setEndOfToken(t)

	return *t, nil
}

// @todo: test all the error cases
func (l *Lexer) eatInputDueToBlockComment() error {
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

		if c == '/' {
			l.eatUntilNewLine()
			return nil
		}
		if c == '*' {
			err = l.eatUntilCharacterCombo('*', '/')
			if err != nil {
				return fmt.Errorf("expected end of block comment, but got end of file: [%s] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
			}

			return nil
		}

		return fmt.Errorf("expected start of block comment, but got unexpected character: [%s] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
	}

	return fmt.Errorf("expected start of block comment, but got unexpected character: [%s] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
}

func (l *Lexer) eatUntilCharacterCombo(r1, r2 rune) error {
	for {
		c, err := l.peekNextCharacter()
		if err != nil {
			return err
		}

		if c == r1 {
			l.eatCharacter()
			c, err = l.peekNextCharacter()
			if err != nil {
				return err
			}

			if c == r2 {
				l.eatCharacter()
				return nil
			}
		}

		l.eatCharacter()
	}
}

func (l *Lexer) eatUntilNewLine() {
	for {
		c, err := l.peekNextCharacter()
		if err != nil {
			break
		}

		if c == '\n' {
			l.eatCharacter()
			return
		}

		l.eatCharacter()
	}
}

func (l *Lexer) composeNewToken() (token.Token, error) {
	t := &token.Token{}
	t.Kind = token.EndOfLine

	for {
		c, err := l.peekNextCharacter()
		if err != nil {
			break
		}

		for token.IsWhitespace(c) {
			l.eatCharacter()
			c, err = l.peekNextCharacter()
			if err != nil {
				return *t, nil
			}
		}

		if c == '/' {
			err = l.eatInputDueToBlockComment()
			if err != nil {
				return *t, err
			}

			continue
		}

		if c == ':' {
			l.eatCharacter()
			c, err = l.peekNextCharacterIgnoringRunes(' ')
			if err != nil {
				return *t, err
			}

			if token.IsStartOfIdentifier(c) {
				t, err := l.makeIdentifierOrKeyword()
				if err != nil {
					return t, err
				}

				t.Kind = token.Label
				return t, nil
			}
		}

		if token.IsStartOfIdentifier(c) {
			return l.makeIdentifierOrKeyword()
		}

		if token.IsPartOfNumber(c) {
			return l.makeNumber()
		} else {
			return *t, fmt.Errorf("unexpected character: [%v] (%d:%d)", string(c), l.CurrentLineNumber, l.CurrentLineCharacterIndex)
		}
	}

	return *t, nil
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
