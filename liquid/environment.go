package liquid

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/astcopy"
	"github.com/spaspa/gorefinement/refinement"
	"golang.org/x/tools/go/analysis"
)

type ObjectRefinementMap = map[types.Object]types.Type

const PredicateVariableName = "__val"

type Environment struct {
	// ExplicitRefinementMap is explicitly defined refinements.
	ExplicitRefinementMap ObjectRefinementMap

	// ImplicitRefinementMap is implicitly reasoned refinements, like assignment.
	ImplicitRefinementMap ObjectRefinementMap

	// FunArgRefinementMap is used to store refinements of arguments passed to function.
	FunArgRefinementMap ObjectRefinementMap

	// Scope is current scope to type check
	Scope *types.Scope

	// Scope is current pos to type check
	Pos token.Pos

	// Pass is current analysis pass
	Pass *analysis.Pass
}

// NewEnvironment creates new type environment with current analysis pass.
func NewEnvironment(pass *analysis.Pass) *Environment {
	return &Environment{
		//TODO make it private and reject object named __val
		ExplicitRefinementMap: map[types.Object]types.Type{},
		ImplicitRefinementMap: map[types.Object]types.Type{},
		FunArgRefinementMap:   map[types.Object]types.Type{},
		Scope:                 nil,
		Pos:                   token.NoPos,
		Pass:                  pass,
	}
}

func (env *Environment) RefinementTypeOf(object types.Object) types.Type {
	e := env.ExplicitRefinementMap[object]
	if e != nil {
		return e
	}
	i := env.ImplicitRefinementMap[object]
	if i != nil {
		return i
	}
	return nil
}

func (env *Environment) Embedding() ast.Expr {
	var result []ast.Expr
	for obj, typ := range env.ExplicitRefinementMap {
		refType, ok := typ.(*refinement.RefinedType)
		if !ok {
			continue
		}

		// refType.Predicate is expr, so this assertion should be safe
		newPredicate := astcopy.Expr(refType.Predicate)
		replaced := replaceIdentOf(newPredicate, refType.RefVar.Name, ast.NewIdent(obj.Name())).(ast.Expr)
		result = append(result, replaced)
	}

	return JoinExpr(token.LAND, result...)
}
