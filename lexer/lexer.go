package lexer

import (
	"log"
	"os"
	"unicode"

	"github.com/fatih/color"
)

type Lexer struct {
	Source []rune
	Cursor int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		Source: []rune(source),
		Cursor: 0,
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
		Source: []rune(string(content)),
		Cursor: 0,
	}
}

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

func PrintDebugTokens(tokens []Token) {
	log.Println("Tokens:")
	for i, token := range tokens {
		log.Printf("[%d] ", i)
		log.Printf("	Token: [%s]\n", token.Kind.String())
		log.Printf("	Value: [%s]\n", token.Value)
		log.Printf("	Location: [%d]\n", token.Location)
	}
}
