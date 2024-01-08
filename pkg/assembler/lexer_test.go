package assembler

import (
	"testing"
)

type Test struct {
	expectedType    int
	expectedLiteral string
}

func TestSimpleLine(t *testing.T) {
	input := `
	push 2067
	`

	tests := []Test{
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "push"},
		{TokenTypeNumber, "2067"},
		{TokenTypeEOL, ""},
	}

	l := NewLexer(input)
	validateTokens(t, tests, l)
}

func TestDataLine(t *testing.T) {
	input := `
	a_number1: w 		1
	`
	tests := []Test{
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "a_number1:"},
		{TokenTypeIdentifier, "w"},
		{TokenTypeNumber, "1"},
		{TokenTypeEOL, ""},
	}

	l := NewLexer(input)
	validateTokens(t, tests, l)
}

func TestSingleLineComments(t *testing.T) {
	input := `
	; This is a header comment
	.Data
		a_number1: 		w		13

	; Start of code
	.CODE
						push	13
						pop

						out ; Ready to write it out!
	`

	tests := []Test{
		{TokenTypeEOL, ""},
		{TokenTypeComment, "; This is a header comment"},
		{TokenTypeEOL, ""},
		{TokenTypeDirective, ".Data"},
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "a_number1:"},
		{TokenTypeIdentifier, "w"},
		{TokenTypeNumber, "13"},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeComment, "; Start of code"},
		{TokenTypeEOL, ""},
		{TokenTypeDirective, ".CODE"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "push"},
		{TokenTypeNumber, "13"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "pop"},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "out"},
		{TokenTypeComment, "; Ready to write it out!"},
		{TokenTypeEOL, ""},
		{TokenTypeEOF, ""},
	}

	l := NewLexer(input)
	validateTokens(t, tests, l)
}

func TestMultilineComments(t *testing.T) {
	input := `
	; This is a header comment
	.Data
		a_number1: 		w		13

	; Start of code
	/*.CODE
						push	13
						pop */

						out ; Ready to write it out!
	`

	tests := []Test{
		{TokenTypeEOL, ""},
		{TokenTypeComment, "; This is a header comment"},
		{TokenTypeEOL, ""},
		{TokenTypeDirective, ".Data"},
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "a_number1:"},
		{TokenTypeIdentifier, "w"},
		{TokenTypeNumber, "13"},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeComment, "; Start of code"},
		{TokenTypeEOL, ""},
		{TokenTypeComment, `/*.CODE
						push	13
						pop */`},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "out"},
		{TokenTypeComment, "; Ready to write it out!"},
		{TokenTypeEOL, ""},
		{TokenTypeEOF, ""},
	}

	l := NewLexer(input)
	validateTokens(t, tests, l)
}

func TestNextToken(t *testing.T) {
	input := `
.DATA
a_number1: w 1
b_number2: w 2.0

.CODE
    ld a_number1
    push 2067
    add
    jmp out
    
out:
    out
`

	tests := []Test{
		{TokenTypeEOL, ""},
		{TokenTypeDirective, ".DATA"},
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "a_number1:"},
		{TokenTypeIdentifier, "w"},
		{TokenTypeNumber, "1"},
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "b_number2:"},
		{TokenTypeIdentifier, "w"},
		{TokenTypeNumber, "2.0"},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeDirective, ".CODE"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "ld"},
		{TokenTypeIdentifier, "a_number1"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "push"},
		{TokenTypeNumber, "2067"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "add"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "jmp"},
		{TokenTypeIdentifier, "out"},
		{TokenTypeEOL, ""},
		{TokenTypeEOL, ""},
		{TokenTypeLabel, "out:"},
		{TokenTypeEOL, ""},
		{TokenTypeIdentifier, "out"},
		{TokenTypeEOL, ""},
		{TokenTypeEOF, ""},
	}

	l := NewLexer(input)
	validateTokens(t, tests, l)
}

func validateTokens(t *testing.T, tests []Test, l *Lexer) {
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s [%s], got=%s [%s]",
				i, getTokenTypeName(tt.expectedType), tt.expectedLiteral, getTokenTypeName(tok.Type), tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
