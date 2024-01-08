package assembler

import (
	"strings"
)

// Token Types
const (
	TokenTypeLabel = iota
	TokenTypeIdentifier
	TokenTypeNumber
	TokenTypeDirective
	TokenTypeComment
	TokenTypeEOL
	TokenTypeEOF
)

var TokenTypeNames []string = []string{
	"label",
	"identifier",
	"number",
	"directive",
	"comment",
	"EOL",
	"EOF",
}

func getTokenTypeName(i int) string {
	if i < 0 || i >= len(TokenTypeNames) {
		return "unknown"
	}

	return TokenTypeNames[i]
}

// Token represents a lexical token.
type Token struct {
	Type    int
	Literal string
	LineNo  int
}

// Lexer represents a lexer.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	prevToken    *Token
	currentLine  int
}

// NewLexer returns a new instance of Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, prevToken: nil, currentLine: 1}
	l.readChar()
	return l
}

// readChar gives us the next character and advance our position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar looks at the next character without moving the current position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) PushToken(t *Token) {
	l.prevToken = t
}

// NextToken returns the next token.
func (l *Lexer) NextToken() *Token {
	if l.prevToken != nil {
		t := l.prevToken
		l.prevToken = nil
		return t
	}

	var tok Token

	l.skipWhitespace()

	switch {
	case isLetter(l.ch):
		literal := l.readIdentifier()
		if strings.HasSuffix(literal, ":") {
			tok = newToken(TokenTypeLabel, literal, l.currentLine)
		} else if strings.HasPrefix(literal, ".") {
			tok = newToken(TokenTypeDirective, literal, l.currentLine)
		} else {
			tok = newToken(TokenTypeIdentifier, literal, l.currentLine)
		}
	case isDigit(l.ch):
		tok = newToken(TokenTypeNumber, l.readNumber(), l.currentLine)
	case l.ch == '\n':
		tok = newToken(TokenTypeEOL, "", l.currentLine)
		l.currentLine++
		l.readChar()
	case l.ch == '.':
		tok = newToken(TokenTypeDirective, l.readIdentifier(), l.currentLine)
	case l.ch == ';': // line comment
		tok = newToken(TokenTypeComment, l.readLineComment(), l.currentLine)
	case l.ch == '/' && l.peekChar() == '*':
		tok = newToken(TokenTypeComment, l.readBlockComment(), l.currentLine)
	case l.ch == 0:
		tok.Literal = ""
		tok.Type = TokenTypeEOF
	default:
		tok = newToken(TokenTypeIdentifier, string(l.ch), l.currentLine)
		l.readChar()
	}

	return &tok
}

// Helper functions

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '.' || l.ch == ':' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readLineComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBlockComment() string {
	position := l.position
	for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
		l.readChar()
	}
	// Move past the closing '*/'
	l.readChar()
	l.readChar()
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType int, literal string, lineNo int) Token {
	return Token{Type: tokenType, Literal: literal, LineNo: lineNo}
}
