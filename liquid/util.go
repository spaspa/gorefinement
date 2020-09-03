package liquid

import (
	"fmt"
	"github.com/spaspa/gorefinement/refinement"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/types"
)

// RefinedTypeFromValue returns narrowest type of given value.
func RefinedTypeFromValue(v constant.Value) (*refinement.RefinedType, error) {
	predicate, err := parser.ParseExpr(fmt.Sprintf("__val == %v", v.ExactString()))
	if err != nil {
		return nil, err
	}
	ident := ast.NewIdent("__val")
	return &refinement.RefinedType{
		Refinement: &refinement.Refinement{
			Predicate: predicate,
			RefVar:    ident,
		},
		Type:       types.Typ[types.UntypedInt],
	}, nil
}
