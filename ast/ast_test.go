package ast

import (
	"testing"

	"github.com/Jitesh117/brainrotLang-interpreter/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&YeetStatement{
				Token: token.Token{Type: token.LET, Literal: "yeet"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "yeet myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
