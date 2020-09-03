package gorefinement

import (
	"fmt"
	"github.com/gostaticanalysis/comment"
	"github.com/spaspa/gorefinement/liquid"
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

	env := liquid.NewEnvironment()

	// TODO: extract type alias

	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil), (*ast.AssignStmt)(nil)}, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			// extract function from FuncDecl
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
			env.ExplicitRefinementMap[pass.TypesInfo.ObjectOf(n.Name)] = refType
		case *ast.AssignStmt:
			// extract explicit variable annotation from definition
			if n.Tok != token.DEFINE {
				return
			}
			lhs := n.Lhs
			rhs := n.Lhs

			// TODO: add support for multiple definition

			commentGroup := cmap.Comments(n)
			if commentGroup == nil {
				// no related comment
				return
			}
			cmt := strings.Replace(commentGroup[0].Text(), "\n", " ", -1)

			if len(lhs) != 1 || len(rhs) != 1 {
				pass.Reportf(n.Pos(), "multiple definition is not supported")
				return
			}

			lhIdent, ok := lhs[0].(*ast.Ident)
			if !ok {
				return
			}

			identStr, refType, err := refinement.ParseWithBaseType(cmt, pass.TypesInfo.TypeOf(lhIdent))
			if err != nil {
				return
			}
			if identStr != lhIdent.Name {
				pass.Reportf(n.Pos(), "[WARN] name of refinement and base did not match: %s, %s", identStr, lhIdent)
				return
			}
			env.ExplicitRefinementMap[pass.TypesInfo.ObjectOf(lhIdent)] = refType
		}
	})

	inspect.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)
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
				typ := argTypeAndValue.Type
				val := argTypeAndValue.Value
				if val != nil {
					if r, err := liquid.RefinedTypeFromValue(val); err == nil {
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
	})

	fmt.Println(env.ExplicitRefinementMap)



	return nil, nil
}
