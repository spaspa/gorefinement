package liquid

import (
	"go/ast"
)

type Result int

const (
	Valid Result = iota
	Invalid
	Unknown
)

// make query ⟦e⟧ ∧ ⟦r1⟧ ⇒ ⟦r2⟧, and validate it.
func Query(env *Environment, p1, p2 ast.Expr) Result {
	return z3Query(env, p1, p2)
}
