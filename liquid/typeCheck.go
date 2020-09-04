package liquid

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/spaspa/gorefinement/freshname"
	"github.com/spaspa/gorefinement/refinement"
)

func TypeCheckExpr(env *Environment, expr ast.Expr) types.Type {
	typeAndValue := env.pass.TypesInfo.Types[expr]
	if typeAndValue.Value != nil {
		if typ, err := refinement.NewRefinedTypeFromValue(typeAndValue.Value); err == nil {
			return typ
		}
	}
	switch expr := expr.(type) {
	case *ast.Ident:
		return checkIdent(env, expr)
	case *ast.CallExpr:
		return checkCallExpr(env, expr)
	case *ast.BinaryExpr:
		return checkOpExpr(env, expr)
	case *ast.UnaryExpr:
		return checkOpExpr(env, expr)
	default:
		return env.pass.TypesInfo.TypeOf(expr)
	}
}

func checkIdent(env *Environment, expr *ast.Ident) types.Type {
	obj := env.pass.TypesInfo.ObjectOf(expr)
	if refType := env.RefinementTypeOf(obj); refType != nil {
		return refType
	} else {
		return env.pass.TypesInfo.TypeOf(expr)
	}
}

func checkCallExpr(env *Environment, expr *ast.CallExpr) types.Type {
	switch fun := expr.Fun.(type) {
	case *ast.Ident:
		obj := env.pass.TypesInfo.ObjectOf(fun)
		sig := env.RefinementTypeOf(obj).(*refinement.DependentSignature)
		return sig.ResultsRefinements
	default:
		return env.pass.TypesInfo.TypeOf(expr)
	}
}

func checkOpExpr(env *Environment, expr ast.Expr) types.Type {
	ident := ast.NewIdent(freshname.Generate())
	return &refinement.RefinedType{
		Refinement: &refinement.Refinement{
			Predicate: &ast.BinaryExpr{
				X:     ident,
				OpPos: token.NoPos,
				Op:    token.EQL,
				Y:     expr,
			},
			RefVar: ident,
		},
		Type: env.pass.TypesInfo.TypeOf(expr),
	}
}
