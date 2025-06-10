package lexer

import (
	"github.com/OG-Open-Source/SDCL/token"
	"unicode/utf8"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '(':
		if l.peekChar() == '(' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DOUBLE_LPAREN, Literal: literal}
		} else {
			tok = newToken(token.LPAREN, l.ch)
		}
	case ')':
		if l.peekChar() == ')' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DOUBLE_RPAREN, Literal: literal}
		} else {
			tok = newToken(token.RPAREN, l.ch)
		}
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
	case '\r':
		if l.peekChar() == '\n' {
			l.readChar() // consume the \n in \r\n
		}
		tok = newToken(token.NEWLINE, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case '-':
		if l.peekChar() == '-' && l.peekChar() == '-' {
			l.readChar()
			l.readChar()
			tok = token.Token{Type: token.FRONTMATTER_DELIMITER, Literal: "---"}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readComment()
		// We don't call l.readChar() here because readComment() does it.
		return tok
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (token.TokenType, string) {
	position := l.position
	var tokType token.TokenType = token.INT
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		tokType = token.FLOAT
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return tokType, l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position + 1 // +1 to skip the '#'
	for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
