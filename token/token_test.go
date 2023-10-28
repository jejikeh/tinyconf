package token

import (
	"testing"

	"github.com/jejikeh/ambient/common"
	"github.com/stretchr/testify/assert"
)

func TestDetectMyKind(t *testing.T) {
	// Test case 1: Testing when value is a keyword

	for name, kind := range keywords {
		t1 := &Token{TokenValue: TokenValue{IndentValue: IndentValue{Name: name, Hash: common.Hash(name)}}}
		t1.DetectMyKind()
		if t1.Kind != kind {
			t.Errorf("Expected t1.Kind to be %v, got %v", kind, t1.Kind)
		}

		t2 := &Token{TokenValue: TokenValue{IndentValue: IndentValue{Name: string(kind), Hash: common.Hash(string(kind))}}}
		t2.DetectMyKind()
		if t2.Kind != Identifier {
			t.Errorf("Expected t2.Kind to be Identifier, got %v", t2.Kind)
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		input rune
		want  bool
	}{
		{' ', true},            // testing space character
		{'\t', true},           // testing tab character
		{'\n', true},           // testing newline character
		{'\r', true},           // testing carriage return character
		{'\v', true},           // testing vertical tab character
		{'\f', true},           // testing form feed character
		{'a', false},           // testing non-whitespace character
		{rune(0x2003), true},   // testing a specific Unicode whitespace character
		{rune(0x1F4A9), false}, // testing a non-whitespace Unicode character
	}

	for _, test := range tests {
		got := IsWhitespace(test.input)
		if got != test.want {
			t.Errorf("IsWhitespace(%v) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestIsStartOfIdentifier(t *testing.T) {
	testData := map[rune]bool{
		'a':  true,  // Testing for a letter
		'_':  true,  // Testing for an underscore
		'5':  false, // Testing for a digit
		'@':  false, // Testing for a special character
		'z':  true,  // Testing for a lowercase letter
		'A':  true,  // Testing for an uppercase letter
		' ':  false, // Testing for a whitespace character
		'\t': false, // Testing for a tab character
	}

	for input, expected := range testData {
		result := IsStartOfIdentifier(input)
		if result != expected {
			t.Errorf("IsStartOfIdentifier(%q) should be %v, got %v", input, expected, result)
		}
	}
}

func TestIsPartOfIdentifier(t *testing.T) {
	tests := []struct {
		name string
		c    rune
		want bool
	}{
		{
			name: "letter",
			c:    'a',
			want: true,
		},
		{
			name: "digit",
			c:    '1',
			want: true,
		},
		{
			name: "underscore",
			c:    '_',
			want: true,
		},
		{
			name: "not identifier",
			c:    '$',
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPartOfIdentifier(tt.c); got != tt.want {
				t.Errorf("IsPartOfIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_SetIndentValue(t *testing.T) {
	t.Run("when s is not empty", func(t *testing.T) {
		token := &Token{}
		token.SetIndentValue("test")

		assert.Equal(t, "test", token.IndentValue.Name)
		assert.Equal(t, common.Hash("test"), token.IndentValue.Hash)
	})
}
