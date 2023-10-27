package lexer

import (
	"fmt"
	"testing"

	"github.com/jejikeh/ambient/common"
	"github.com/jejikeh/ambient/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveLabelIdentifierDeclaration(t *testing.T) {
	lexer := Lexer{}

	// Test case 1: No label tokens
	tokens := []token.Token{
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "x"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "y"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "z"}}},
	}
	expected := []token.Token{
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "x"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "y"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "z"}}},
	}
	result := lexer.resolveLabelIdentifierDeclaration(tokens)
	assert.Equal(t, expected, result, "Failed test case 1")

	// Test case 2: With label tokens
	tokens = []token.Token{
		{Kind: token.Label, LineStart: 12, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label1"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "y"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label1"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "z"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "x"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label2"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "h"}}},
	}
	expected = []token.Token{
		{Kind: token.Label, LineStart: 12, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label1"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "y"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label1"}, IntegerValue: 0}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "z"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "x"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "label2"}}},
		{Kind: token.Identifier, TokenValue: token.TokenValue{IndentValue: token.IndentValue{Name: "h"}}},
	}
	result = lexer.resolveLabelIdentifierDeclaration(tokens)
	assert.Equal(t, expected, result, "Failed test case 2")
}

func TestLexer_peekNextCharacterIgnoringRunes(t *testing.T) {
	l := NewLexer("abc")

	// Test case: rune is not present
	r, err := l.peekNextCharacterIgnoringRunes('d')
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if r != 'a' {
		t.Errorf("Expected 'a', got %c", r)
	}

	// Test case: rune is present
	r, err = l.peekNextCharacterIgnoringRunes('b')
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if r != 'a' {
		t.Errorf("Expected 'a', got %c", r)
	}

	// Test case: multiple runes present
	r, err = l.peekNextCharacterIgnoringRunes('a')
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if r != 'b' {
		t.Errorf("Expected 'b', got %c", r)
	}
}

func TestLexer_parseIndentifier(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid identifier",
			input:   "abc123",
			wantErr: false,
		},
		{
			name:    "Invalid identifier",
			input:   "123abc",
			wantErr: true,
		},
		{
			name:    "Empty identifier",
			input:   "",
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := &token.Token{}

			err := l.parseIndentifier(tok)

			if (err != nil) != tt.wantErr {
				t.Errorf("[%d] expected error: %v, got: %v", i, tt.wantErr, err)
			}
		})
	}
}

func TestLexer_makeIdentifierOrKeyword(t *testing.T) {
	lexer := NewLexer("loop_label")

	// Testing identifier
	t.Run("Identifier", func(t *testing.T) {
		expectedToken := token.Token{
			Kind: token.Identifier,
			TokenValue: token.TokenValue{
				IndentValue: token.IndentValue{
					Name: "loop_label",
					Hash: common.Hash("loop_label"),
				},
			},
			LineStart:    0,
			LineEnd:      0,
			CollumnStart: 0,
			CollumnEnd:   10,
		}

		actualToken, err := lexer.makeIdentifierOrKeyword()
		require.NoError(t, err)
		assert.Equal(t, expectedToken, actualToken)
	})

	lexer1 := NewLexer("jmp")

	// Testing identifier
	t.Run("Identifier", func(t *testing.T) {
		expectedToken := token.Token{
			Kind: token.Jump,
			TokenValue: token.TokenValue{
				IndentValue: token.IndentValue{
					Name: "jmp",
					Hash: common.Hash("jmp"),
				},
			},
			LineStart:    0,
			LineEnd:      0,
			CollumnStart: 0,
			CollumnEnd:   3,
		}

		actualToken, err := lexer1.makeIdentifierOrKeyword()
		require.NoError(t, err)
		assert.Equal(t, expectedToken, actualToken)
	})
}

func TestLexer_makeNumber(t *testing.T) {
	// Test case 1: Number starts with a digit
	t.Run("Number starts with a digit", func(t *testing.T) {
		l := NewLexer("123")
		expected := token.Token{
			Kind: token.Number,
			TokenValue: token.TokenValue{
				IndentValue: token.IndentValue{
					Name: "123",
					Hash: common.Hash("123"),
				},
				IntegerValue: 123,
			},
			LineStart:    0,
			LineEnd:      0,
			CollumnStart: 0,
			CollumnEnd:   3,
		}

		actual, err := l.makeNumber()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("Expected token %v, but got %v", expected, actual)
		}
	})

	// Test case 2: Number starts with a non-digit character
	t.Run("Number starts with a non-digit character", func(t *testing.T) {
		l := NewLexer("#123")

		expectedErr := fmt.Errorf("expected start of number, but got unexpected character: [#] (0:0)")

		_, err := l.makeNumber()
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, but got %v", expectedErr, err)
		}
	})

	// Test case 3: Number is empty
	t.Run("Number is empty", func(t *testing.T) {
		l := NewLexer("")
		expectedErr := fmt.Errorf(`unexpected end of file!
		Current cursor: [0]
		Source length: [0]
		Current line: [0]
		Current column: [0]`)
		_, err := l.makeNumber()
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, but got %v", expectedErr, err)
		}
	})
}

func TestLexer_eatInputDueToBlockComment(t *testing.T) {
	lexer := NewLexer("/* test */ hello /* world */")
	err := lexer.eatInputDueToBlockComment()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lexer = NewLexer("// test")
	err = lexer.eatInputDueToBlockComment()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lexer = NewLexer("// test */")
	err = lexer.eatInputDueToBlockComment()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	lexer = NewLexer("test")
	err = lexer.eatInputDueToBlockComment()
	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	lexer = NewLexer("/* test")
	err = lexer.eatInputDueToBlockComment()
	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	lexer = NewLexer("/* test //")
	err = lexer.eatInputDueToBlockComment()
	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	lexer = NewLexer("test */")
	err = lexer.eatInputDueToBlockComment()
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
}

func TestLexer_eatUntilCharacterCombo(t *testing.T) {

	lexer := NewLexer("helloworld")
	// Test case 1: Testing when r1 and r2 are both present in the input string
	err := lexer.eatUntilCharacterCombo('o', 'w')
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test case 2: Testing when r1 and r2 are not present in the input string
	lexer = NewLexer("hello")
	err = lexer.eatUntilCharacterCombo('o', 'w')
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}

	// Test case 3: Testing when r1 is present but r2 is not present in the input string
	lexer = NewLexer("hello world")
	err = lexer.eatUntilCharacterCombo('o', 'x')
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
}

func TestLexer_composeNewToken(t *testing.T) {
	l := &Lexer{}

	testCases := map[string]struct {
		input       string
		expectedTok token.Token
		expectedErr error
	}{
		"Whitespace": {
			input:       "   ",
			expectedTok: token.Token{Kind: token.EndOfLine},
		},
		"BlockComment": {
			input:       "// Block comment\n",
			expectedTok: token.Token{Kind: token.EndOfLine},
		},
		"Label": {
			input:       ": label",
			expectedTok: token.Token{Kind: token.Label},
		},
		"IdentifierOrKeyword": {
			input:       "identifier",
			expectedTok: token.Token{Kind: token.Identifier},
		},
		"Number": {
			input:       "123",
			expectedTok: token.Token{Kind: token.Number},
		},
		"UnexpectedCharacter": {
			input:       "%",
			expectedTok: token.Token{Kind: token.EndOfLine},
			expectedErr: fmt.Errorf("unexpected character: [%v] (0:0)", "%"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			l = NewLexer(tc.input)

			tok, err := l.composeNewToken()
			if err != nil && tc.expectedErr == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && tc.expectedErr != nil {
				t.Errorf("Expected error: %v, got nil", tc.expectedErr)
			}
			if tok.Kind != tc.expectedTok.Kind {
				t.Errorf("Expected token kind %v, got %v", tc.expectedTok.Kind, tok.Kind)
			}
		})
	}
}
