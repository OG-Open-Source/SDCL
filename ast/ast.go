package ast

import "github.com/OG-Open-Source/SDCL/token"

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// Statement is a node that represents a statement.
type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type ObjectLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (ol *ObjectLiteral) expressionNode()      {}
func (ol *ObjectLiteral) TokenLiteral() string { return ol.Token.Literal }

// Expression is a node that represents an expression.
type Expression interface {
	Node
	expressionNode()
}

// Document is the root node of any SDCL AST.
type Document struct {
	Statements []Statement
}

func (d *Document) TokenLiteral() string {
	if len(d.Statements) > 0 {
		return d.Statements[0].TokenLiteral()
	}
	return ""
}

// KeyValuePair represents a key-value pair.
type KeyValuePair struct {
	Token token.Token // the KEY token
	Key   *Identifier
	Value Expression
}

func (kv *KeyValuePair) statementNode()       {}
func (kv *KeyValuePair) TokenLiteral() string { return kv.Token.Literal }

// Identifier represents a key.
type Identifier struct {
	Token token.Token // the KEY token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Literal represents a literal value.
type Literal interface {
	Expression
	literalNode()
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// FloatLiteral represents a float literal.
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// BooleanLiteral represents a boolean literal.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }

// NullLiteral represents a null literal.
type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Literal }

// DateLiteral represents a date literal.
type DateLiteral struct {
	Token token.Token
	Value string
}

func (dl *DateLiteral) expressionNode()      {}
func (dl *DateLiteral) TokenLiteral() string { return dl.Token.Literal }

// TimeLiteral represents a time literal.
type TimeLiteral struct {
	Token token.Token
	Value string
}

func (tl *TimeLiteral) expressionNode()      {}
func (tl *TimeLiteral) TokenLiteral() string { return tl.Token.Literal }

// DateTimeLiteral represents a datetime literal.
type DateTimeLiteral struct {
	Token token.Token
	Value string
}

func (dtl *DateTimeLiteral) expressionNode()      {}
func (dtl *DateTimeLiteral) TokenLiteral() string { return dtl.Token.Literal }

// CountryLiteral represents a country code literal.
type CountryLiteral struct {
	Token token.Token
	Value string
}

func (cl *CountryLiteral) expressionNode()      {}
func (cl *CountryLiteral) TokenLiteral() string { return cl.Token.Literal }

// Base64Literal represents a base64 literal.
type Base64Literal struct {
	Token token.Token
	Value string
}

func (b64l *Base64Literal) expressionNode()      {}
func (b64l *Base64Literal) TokenLiteral() string { return b64l.Token.Literal }

// ObjectLiteral represents an object literal.

// ArrayLiteral represents an array literal.
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// ValueReference represents a value reference.
type ValueReference struct {
	Token token.Token // the '(' token
	Path  []*Identifier
}

func (vr *ValueReference) expressionNode()      {}
func (vr *ValueReference) TokenLiteral() string { return vr.Token.Literal }

// ContentInclusion represents a content inclusion.
type ContentInclusion struct {
	Token token.Token // the '((' token
	Path  []*Identifier
}

func (ci *ContentInclusion) expressionNode()      {}
func (ci *ContentInclusion) TokenLiteral() string { return ci.Token.Literal }

// ExternalReference represents an external reference.
type ExternalReference struct {
	Token token.Token // the '.' token
	Path  []*Identifier
}

func (er *ExternalReference) expressionNode()      {}
func (er *ExternalReference) TokenLiteral() string { return er.Token.Literal }