package liquid

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"github.com/aclements/go-z3/z3"
)

// make query ⟦e⟧ ∧ ⟦r1⟧ ⇒ ⟦r2⟧, and validate it.
func z3Query(env *Environment, p1, p2 ast.Expr) Result {
	// context
	config := z3.NewContextConfig()
	ctx := z3.NewContext(config)

	envEmbedding := env.Embedding()
	antecedent, err := convertToZ3Ast(env, ctx, JoinExpr(token.LAND, envEmbedding, p1))
	if err != nil {
		log.Println("antecedent expr construction failed:", err)
		return Unknown
	}
	consequent, err := convertToZ3Ast(env, ctx, p2)
	if err != nil {
		log.Println("consequent expr construction failed:", err)
		return Unknown
	}

	antecedentBool, ok := antecedent.(z3.Bool)
	if !ok {
		log.Println("antecedent is not bool")
		return Unknown
	}
	consequentBool, ok := consequent.(z3.Bool)
	if !ok {
		log.Println("consequent is not bool")
		return Unknown
	}

	// solver
	s := z3.NewSolver(ctx)

	z3Ast := antecedentBool.Implies(consequentBool).Not()
	fmt.Println(z3Ast)
	s.Assert(z3Ast)

	sat, err := s.Check()

	if err != nil {
		log.Println("sat check failed:", err)
		return Unknown
	}
	if sat {
		return Invalid
	} else {
		return Valid
	}
}
