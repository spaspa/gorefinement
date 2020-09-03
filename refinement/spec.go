package refinement

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/types"
	"strings"
)

var truePredicate, _ = parser.ParseExpr("true")

type RefinedType struct {
	*Refinement
	types.Type
	isConstant bool
}

// NewRefinedTypeFromValue returns narrowest type of given value.
func NewRefinedTypeFromValue(v constant.Value) (*RefinedType, error) {
	predicate, err := parser.ParseExpr(fmt.Sprintf("__val == %v", v.ExactString()))
	if err != nil {
		return nil, err
	}
	ident := ast.NewIdent("__val")
	return &RefinedType{
		Refinement: &Refinement{
			Predicate: predicate,
			RefVar:    ident,
		},
		Type:       types.Typ[types.UntypedInt],
		isConstant: true,
	}, nil
}
func NewRefinedTypeWithTruePredicate(typ types.Type) *RefinedType {
	return &RefinedType{
		Refinement: &Refinement{
			Predicate: truePredicate,
			RefVar:    ast.NewIdent("_"),
		},
		Type: typ,
	}
}
func (r *RefinedType) Underlying() types.Type { return r }
func (r *RefinedType) String() string {
	if r == nil {
		return "nil"
	}
	return fmt.Sprintf("{ %s: %s | %s }", r.RefVar.Name, r.Type, types.ExprString(r.Predicate))
}
func (r *RefinedType) IsConstant() bool {
	return r.isConstant
}
func (r *RefinedType) ConstantNode() ast.Expr {
	if !r.isConstant {
		return nil
	}
	return r.Predicate.(*ast.BinaryExpr).Y
}

type DependentSignature struct {
	*types.Signature
	ParamRefinements   *RefinedTuple
	ResultsRefinements *RefinedTuple
}

func NewDependentSignature(sig *types.Signature, params, results *RefinedTuple) *DependentSignature {
	return &DependentSignature{
		sig,
		params,
		results,
	}
}
func (s *DependentSignature) Underlying() types.Type { return s }
func (s *DependentSignature) String() string {
	if s == nil {
		return "nil"
	}
	var paramStr []string
	var resultStr []string

	for i := 0; i < s.ParamRefinements.Len(); i++ {
		param := s.ParamRefinements.At(i)
		paramStr = append(paramStr, fmt.Sprintf("%s: %s", param.name, param.RefinedType))
	}
	for i := 0; i < s.ResultsRefinements.Len(); i++ {
		result := s.ResultsRefinements.At(i)
		if result.name == "" {
			resultStr = append(resultStr, fmt.Sprintf("%s", result.RefinedType))
		} else {
			resultStr = append(resultStr, fmt.Sprintf("%s: %s", result.name, result.RefinedType))
		}
	}

	if s.ResultsRefinements.Len() == 1 {
		return fmt.Sprintf("(%s) -> %s", strings.Join(paramStr, ", "), strings.Join(resultStr, ", "))
	} else {
		return fmt.Sprintf("(%s) -> (%s)", strings.Join(paramStr, ", "), strings.Join(resultStr, ", "))
	}
}

type Refinement struct {
	Predicate ast.Expr
	RefVar    *ast.Ident
}

func (r *Refinement) String() string {
	return fmt.Sprintf("{ %s: ? | %s }", r.RefVar.Name, types.ExprString(r.Predicate))
}

type RefinedVar struct {
	// TODO
	name        string
	RefinedType *RefinedType
}

type RefinedTuple struct {
	vars []*RefinedVar
}

func NewRefinedTuple(x ...*RefinedVar) *RefinedTuple {
	if len(x) > 0 {
		return &RefinedTuple{x}
	}
	return nil
}
func (t *RefinedTuple) Underlying() types.Type {
	return t
}
func (t *RefinedTuple) String() string {
	return "##refined tuple##"
}
func (t *RefinedTuple) Len() int {
	if t != nil {
		return len(t.vars)
	}
	return 0
}
func (t *RefinedTuple) At(i int) *RefinedVar {
	return t.vars[i]
}
