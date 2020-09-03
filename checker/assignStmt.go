package checker

import (
	"github.com/spaspa/gorefinement/liquid"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

func CheckAssignStmt(pass *analysis.Pass, env *liquid.Environment, assignStmt *ast.AssignStmt) {
	if assignStmt.Tok != token.ASSIGN {
		return
	}
	if len(assignStmt.Lhs) > 1 || len(assignStmt.Rhs) > 1 {
		// TODO: Support multiple assignment
		return
	}

	lhs, _ := assignStmt.Lhs[0].(*ast.Ident)
	if lhs == nil {
		return
	}
	rhs := assignStmt.Rhs[0]

	env.Pos = assignStmt.Pos()
	env.Scope = pass.Pkg.Scope().Innermost(assignStmt.Pos())

	lhsType := liquid.TypeCheckExpr(env, lhs)
	rhsType := liquid.TypeCheckExpr(env, rhs)

	result := liquid.IsSubtype(env, rhsType, lhsType)
	if !result {
		pass.Reportf(assignStmt.Pos(), "UNSAFE")
	}
}
