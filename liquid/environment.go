package liquid

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/astcopy"
	"github.com/spaspa/gorefinement/refinement"
	"golang.org/x/tools/go/analysis"
)

type ObjectRefinementMap map[types.Object]types.Type
type NameRefinementMap map[string]types.Type

const predicateVariableName = "__val"
const argumentVariablePrefix = "__arg_"

type Environment struct {
	// ExplicitRefinementMap is explicitly defined refinements.
	ExplicitRefinementMap ObjectRefinementMap

	// ImplicitRefinementMap is implicitly reasoned refinements, like assignment.
	ImplicitRefinementMap ObjectRefinementMap

	// funArgRefinementMap is used to store refinements of arguments passed to function.
	// Argument labels cannot have a corresponding object, so store type by its name.
	funArgRefinementMap NameRefinementMap

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
		funArgRefinementMap:   map[string]types.Type{},
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

func (env *Environment) AddFunArgRefinement(label string, typ types.Type) {
	env.funArgRefinementMap[argumentVariablePrefix+label] = typ
}
func (env *Environment) ClearFunArgRefinement() {
	env.funArgRefinementMap = map[string]types.Type{}
}

func (env *Environment) Embedding() ast.Expr {
	var result []ast.Expr
	result = append(result, env.ExplicitRefinementMap.collectExpr(env)...)
	result = append(result, env.ExplicitRefinementMap.collectExpr(env)...)
	result = append(result, env.funArgRefinementMap.collectExpr()...)
	return JoinExpr(token.LAND, result...)
}

func (m ObjectRefinementMap) collectExpr(env *Environment) []ast.Expr {
	var result []ast.Expr
	if m == nil {
		return result
	}
	for obj, typ := range m {
		if _, o := env.Scope.LookupParent(obj.Name(), env.Pos); o == nil {
			// not in scope
			continue
		}
		refType, ok := typ.(*refinement.RefinedType)
		if !ok {
			continue
		}

		// refType.Predicate is expr, so this assertion should be safe
		newPredicate := astcopy.Expr(refType.Predicate)
		replaced := replaceIdentOf(newPredicate, refType.RefVar.Name, ast.NewIdent(obj.Name())).(ast.Expr)
		result = append(result, replaced)
	}
	return result
}

func (m NameRefinementMap) collectExpr() []ast.Expr {
	var result []ast.Expr
	if m == nil {
		return result
	}
	for name, typ := range m {
		refType, ok := typ.(*refinement.RefinedType)
		if !ok {
			continue
		}

		// refType.Predicate is expr, so this assertion should be safe
		newPredicate := astcopy.Expr(refType.Predicate)
		replaced := replaceIdentOf(newPredicate, refType.RefVar.Name, ast.NewIdent(name)).(ast.Expr)
		result = append(result, replaced)
	}
	return result
}
