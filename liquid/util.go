package liquid

import (
	"go/ast"
	"go/token"

	"github.com/go-toolsmith/astcopy"
	"github.com/spaspa/gorefinement/refinement"
	"golang.org/x/tools/go/ast/astutil"
)

func replaceIdentOf(n ast.Node, name string, to ast.Node) ast.Node {
	return astutil.Apply(n, func(cursor *astutil.Cursor) bool {
		current := cursor.Node()
		if ident, ok := current.(*ast.Ident); ok && ident.Name == name {
			cursor.Replace(to)
		}
		return true
	}, nil)
}

func joinExpr(sep token.Token, es ...ast.Expr) ast.Expr {
	if len(es) == 0 {
		return nil
	}

	var result = es[0]
	for i := 1; i < len(es); i++ {
		result = &ast.BinaryExpr{
			X:     result,
			OpPos: token.NoPos,
			Op:    sep,
			Y:     es[i],
		}
	}

	return result
}

func normalizedPredicateOf(r *refinement.RefinedType) ast.Expr {
	newPred := astcopy.Expr(r.Predicate)
	return replaceIdentOf(newPred, r.RefVar.Name, ast.NewIdent(predicateVariableName)).(ast.Expr)
}
