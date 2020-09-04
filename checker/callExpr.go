package checker

import (
	"go/ast"

	"github.com/spaspa/gorefinement/liquid"
	"github.com/spaspa/gorefinement/refinement"
	"golang.org/x/tools/go/analysis"
)

func CheckCallExpr(pass *analysis.Pass, env *liquid.Environment, callExpr *ast.CallExpr) {
	funIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return
	}
	funDepSig, ok := liquid.TypeCheckExpr(env, funIdent).(*refinement.DependentSignature)
	if !ok {
		return
	}

	env.Pos = callExpr.Pos()
	env.Scope = pass.Pkg.Scope().Innermost(callExpr.Pos())

	for i, arg := range callExpr.Args {
		typ := liquid.TypeCheckExpr(env, arg)

		argVar := funDepSig.ParamRefinements.At(i)

		result := liquid.IsSubtype(env, typ, funDepSig.ParamRefinements.At(i).RefinedType)
		if !result {
			pass.Reportf(callExpr.Pos(), "UNSAFE")
		}

		env.AddFunArgRefinement(argVar.Name, typ)
	}
	env.ClearFunArgRefinement()
}
