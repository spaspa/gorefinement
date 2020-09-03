package liquid

import (
	"github.com/spaspa/gorefinement/refinement"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

type ObjectRefinementMap = map[types.Object]types.Type

type Environment struct {
	// ExplicitRefinementMap is explicitly defined refinements.
	ExplicitRefinementMap ObjectRefinementMap

	// ImplicitRefinementMap is implicitly reasoned refinements, like assignment.
	ImplicitRefinementMap ObjectRefinementMap

	// FunArgRefinementMap is used to store refinements of arguments passed to function.
	FunArgRefinementMap ObjectRefinementMap

	// Scope is current scope to type check
	Scope *types.Scope

	// analysis pass
	pass *analysis.Pass
}

// NewEnvironment creates new type environment with current analysis pass.
func NewEnvironment(pass *analysis.Pass) *Environment {
	return &Environment{
		ExplicitRefinementMap: map[types.Object]types.Type{},
		ImplicitRefinementMap: map[types.Object]types.Type{},
		FunArgRefinementMap:   map[types.Object]types.Type{},
		Scope:                 nil,
		pass:                  pass,
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
		replaced := ReplaceIdentOf(refType.Predicate, refType.RefVar.Name, ast.NewIdent(obj.Name())).(ast.Expr)
		result = append(result, replaced)
	}

	return JoinExpr(result, token.LAND)
}

