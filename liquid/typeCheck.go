package liquid

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/spaspa/gorefinement/freshname"
	"github.com/spaspa/gorefinement/refinement"
)

func TypeCheckExpr(env *Environment, expr ast.Expr) types.Type {
	typeAndValue := env.Pass.TypesInfo.Types[expr]
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
		return env.Pass.TypesInfo.TypeOf(expr)
	}
}

func checkIdent(env *Environment, expr *ast.Ident) types.Type {
	obj := env.Pass.TypesInfo.ObjectOf(expr)
	if refType := env.RefinementTypeOf(obj); refType != nil {
		return refType
	} else {
		return env.Pass.TypesInfo.TypeOf(expr)
	}
}

func checkCallExpr(env *Environment, expr *ast.CallExpr) types.Type {
	switch fun := expr.Fun.(type) {
	case *ast.Ident:
		obj := env.Pass.TypesInfo.ObjectOf(fun)
		sig := env.RefinementTypeOf(obj).(*refinement.DependentSignature)
		return sig.ResultsRefinements
	default:
		return env.Pass.TypesInfo.TypeOf(expr)
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
		Type: env.Pass.TypesInfo.TypeOf(expr),
	}
}

func checkBinaryExpr(env *Environment, expr *ast.BinaryExpr) types.Type {
	ty1 := TypeCheckExpr(env, expr.X)
	ty2 := TypeCheckExpr(env, expr.Y)
	if ty1 == nil || ty2 == nil {
		return env.Pass.TypesInfo.TypeOf(expr)
	}

	refType1, _ := ty1.(*refinement.RefinedType)
	refType2, _ := ty2.(*refinement.RefinedType)

	if refType1 == nil || refType2 == nil {
		return env.Pass.TypesInfo.TypeOf(expr)
	}

	if refType1.IsConstant() {
		return &refinement.RefinedType{
			Refinement: &refinement.Refinement{
				Predicate: &ast.BinaryExpr{
					X:     refType1.ConstantNode(),
					OpPos: token.NoPos,
					Op:    expr.Op,
					Y:     refType2.RefVar,
				},
				RefVar: refType1.RefVar,
			},
			Type: refType2.Type,
		}
	}
	if refType2.IsConstant() {
		return &refinement.RefinedType{
			Refinement: &refinement.Refinement{
				Predicate: &ast.BinaryExpr{
					X:     refType1.RefVar,
					OpPos: token.NoPos,
					Op:    expr.Op,
					Y:     refType2.ConstantNode(),
				},
				RefVar: refType1.RefVar,
			},
			Type: refType1.Type,
		}
	}

	var baseType types.Type

	if types.AssignableTo(refType1.Type, refType2.Type) {
		baseType = refType1.Type
	} else if types.AssignableTo(refType2.Type, refType1.Type) {
		baseType = refType2.Type
	} else {
		return env.Pass.TypesInfo.TypeOf(expr)
	}

	return &refinement.RefinedType{
		Refinement: &refinement.Refinement{
			Predicate: &ast.BinaryExpr{
				X:     refType1.RefVar,
				OpPos: token.NoPos,
				Op:    expr.Op,
				Y:     refType2.RefVar,
			},
			RefVar: refType1.RefVar,
		},
		Type: baseType,
	}

}

func checkUnaryExpr(env *Environment, expr *ast.UnaryExpr) types.Type {
	ty1 := TypeCheckExpr(env, expr.X)
	if ty1 == nil {
		return env.Pass.TypesInfo.TypeOf(expr)
	}

	refType1, _ := ty1.(*refinement.RefinedType)

	if refType1 == nil {
		return env.Pass.TypesInfo.TypeOf(expr)
	}

	return &refinement.RefinedType{
		Refinement: &refinement.Refinement{
			Predicate: &ast.UnaryExpr{
				OpPos: token.NoPos,
				Op:    expr.Op,
				X:     nil,
			},
			RefVar: refType1.RefVar,
		},
		Type: refType1.Type,
	}
}
