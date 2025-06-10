package parser

import (
	"testing"

	"github.com/OG-Open-Source/SDCL/ast"
	"github.com/OG-Open-Source/SDCL/lexer"
)

func TestParseObjectLiteral(t *testing.T) {
	input := `{ "key": "value", "another.key": 123 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	obj, ok := stmt.Expression.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ObjectLiteral. got=%T", stmt.Expression)
	}
	if len(obj.Pairs) != 2 {
		t.Errorf("obj.Pairs has wrong number of pairs. got=%d", len(obj.Pairs))
	}
}