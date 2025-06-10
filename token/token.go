package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	KEY    = "KEY"    // add, foobar, x, y, ...
	STRING = "STRING" // "foobar"
	INT    = "INT"    // 12345
	FLOAT  = "FLOAT"  // 12.34
	
	// Operators and Delimiters
	COLON      = ":"
	LBRACE     = "{"
	RBRACE     = "}"
	LBRACKET   = "["
	RBRACKET   = "]"
	LPAREN     = "("
	RPAREN     = ")"
	DOUBLE_LPAREN = "(("
	DOUBLE_RPAREN = "))"
	DOT        = "."
	HASH       = "#"
	COMMA      = ","

	// Keywords
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	DATE     = "DATE"
	TIME     = "TIME"
	DATETIME = "DATETIME"
	COUNTRY  = "COUNTRY"
	BASE64   = "BASE64"
	
	// Special
	FRONTMATTER_DELIMITER = "---"
	WHITESPACE = "WHITESPACE"
	NEWLINE = "NEWLINE"
	COMMENT = "COMMENT"
)

var keywords = map[string]TokenType{
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"date":     DATE,
	"time":     TIME,
	"datetime": DATETIME,
	"country":  COUNTRY,
	"base64":   BASE64,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return KEY
}