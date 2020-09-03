package checker

import (
	"go/ast"
	"go/types"

	"github.com/spaspa/gorefinement/liquid"
	"github.com/spaspa/gorefinement/refinement"
	"golang.org/x/tools/go/analysis"
)

func CheckCallExpr(pass *analysis.Pass, env *liquid.Environment, callExpr *ast.CallExpr) {
	funIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok {
		return
	}
	funObj := pass.TypesInfo.ObjectOf(funIdent)
	if funObj == nil {
		return
	}
	funDepSig, _ := env.RefinementTypeOf(funObj).(*refinement.DependentSignature)
	if funDepSig == nil {
		return
	}

	env.Scope = pass.Pkg.Scope().Innermost(callExpr.Pos())

	for i, arg := range callExpr.Args {
		var checkType types.Type
		switch arg := arg.(type) {
		case *ast.Ident:
			argObj := pass.TypesInfo.ObjectOf(arg)
			argRefType := env.RefinementTypeOf(argObj)
			if argRefType != nil {
				checkType = argRefType
			} else {
				checkType = pass.TypesInfo.TypeOf(arg)
			}
		default:
			argTypeAndValue := pass.TypesInfo.Types[arg]
			typ := liquid.TypeCheckExpr(env, arg)
			if typ == nil {
				typ = argTypeAndValue.Type
			}

			val := argTypeAndValue.Value
			if val != nil {
				if r, err := refinement.NewRefinedTypeFromValue(val); err == nil {
					typ = r
				}
			}

			checkType = typ
		}
		result := liquid.IsSubtype(env, checkType, funDepSig.ParamRefinements.At(i).RefinedType)
		if !result {
			pass.Reportf(callExpr.Pos(), "UNSAFE")
		}
	}
}
