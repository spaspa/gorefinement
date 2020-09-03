package liquid

import (
	"fmt"
	"github.com/mitchellh/go-z3"
	"go/token"
	"log"

	"go/ast"
)

type Result int

const (
	Valid Result = iota
	Invalid
	Unknown
)

// make query ⟦e⟧ ∧ ⟦r1⟧ ⇒ ⟦r2⟧, and validate it.
func query(env *Environment, p1, p2 ast.Expr) Result {
	// context
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	err := config.Close()
	if err != nil {
		log.Println("z3 config close failed", err)
	}
	defer ctx.Close()

	envEmbedding := env.Embedding()
	antecedent, err := ConvertToZ3Ast(env, ctx, joinExpr(token.LAND, envEmbedding, p1))
	if err != nil {
		log.Println("antecedent expr construction failed", err)
		return Unknown
	}
	consequent, err := ConvertToZ3Ast(env, ctx, p2)
	if err != nil {
		log.Println("consequent expr construction failed", err)
		return Unknown
	}

	// solver
	s := ctx.NewSolver()
	defer s.Close()

	z3Ast := antecedent.Implies(consequent).Not()
	fmt.Println(z3Ast)
	s.Assert(z3Ast)

	switch s.Check() {
	case z3.True:
		return Invalid
	case z3.False:
		return Valid
	}
	return Unknown
}