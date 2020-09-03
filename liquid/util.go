package liquid

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

func ReplaceIdentOf(n ast.Node, name string, to ast.Node) ast.Node {
	return astutil.Apply(n, func(cursor *astutil.Cursor) bool {
		current := cursor.Node()
		if ident, ok := current.(*ast.Ident); ok && ident.Name == name {
			cursor.Replace(to)
		}
		return true
	}, nil)
}

func JoinExpr(es []ast.Expr, sep token.Token) ast.Expr {
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
