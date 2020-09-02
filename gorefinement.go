package gorefinement

import (
	"fmt"
	"github.com/gostaticanalysis/comment"
	"github.com/spaspa/gorefinement/refinement"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"strings"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "gorefinement is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gorefinement",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cmap := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	explicitRefinementTypesMap := map[types.Object]types.Type{}

	// TODO: extract type alias

	// extract function from FuncDecl
	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			cmt := strings.Replace(n.Doc.Text(), "\n", " ", -1)
			if cmt == "" {
				return
			}
			identStr, refType, err := refinement.ParseWithBaseType(cmt, pass.TypesInfo.TypeOf(n.Name))
			if err != nil {
				return
			}
			if identStr != n.Name.Name {
				pass.Reportf(n.Pos(), "[WARN] name of refinement and base did not match: %s, %s", identStr, n.Name)
				return
			}
			explicitRefinementTypesMap[pass.TypesInfo.ObjectOf(n.Name)] = refType
		}
	})

	// extract explicit variable annotation from definition
	inspect.Preorder([]ast.Node{(*ast.AssignStmt)(nil)}, func(n ast.Node) {
		assignStmt := n.(*ast.AssignStmt)
		if assignStmt.Tok != token.DEFINE {
			return
		}
		lhs := assignStmt.Lhs
		rhs := assignStmt.Lhs

		// TODO: add support for multiple definition

		commentGroup := cmap.Comments(assignStmt)
		if commentGroup == nil {
			// no related comment
			return
		}
		cmt := strings.Replace(commentGroup[0].Text(), "\n", " ", -1)

		if len(lhs) != 1 || len(rhs) != 1 {
			pass.Reportf(assignStmt.Pos(), "multiple definition is not supported")
			return
		}

		lhIdent, ok := lhs[0].(*ast.Ident)
		if !ok {
			return
		}

		identStr, refType, err := refinement.ParseWithBaseType(cmt, pass.TypesInfo.TypeOf(lhIdent))
		fmt.Println(identStr, refType, err)
		explicitRefinementTypesMap[pass.TypesInfo.ObjectOf(lhIdent)] = refType
	})

	fmt.Println(explicitRefinementTypesMap)

	initOrder := pass.TypesInfo.InitOrder
	for _, initializer := range initOrder {
		fmt.Println(initializer)
	}

	return nil, nil
}
