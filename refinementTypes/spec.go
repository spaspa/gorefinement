package refinementTypes

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

type RefinedType struct {
	Refinement
	types.Type
}

func (r *RefinedType) Underlying() types.Type { return r }
func (r *RefinedType) String() string {
	return fmt.Sprintf("{ %s: %s | %s }", r.RefVar.Name, r.Type, types.ExprString(r.Predicate))
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
	var paramStr []string
	var resultStr []string

	for i := 0; i < s.ParamRefinements.Len(); i++ {
		param := s.ParamRefinements.At(i)
		paramStr = append(paramStr, fmt.Sprintf("%s: %s", param.name, param.refinedType))
	}
	for i := 0; i < s.ResultsRefinements.Len(); i++ {
		result := s.ResultsRefinements.At(i)
		if result.name == "" {
			resultStr = append(resultStr, fmt.Sprintf("%s", result.refinedType))
		} else {
			resultStr = append(resultStr, fmt.Sprintf("%s: %s", result.name, result.refinedType))
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

type RefinedVar struct {
	// TODO
	name        string
	refinedType *RefinedType
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
func (t *RefinedTuple) Len() int {
	if t != nil {
		return len(t.vars)
	}
	return 0
}
func (t *RefinedTuple) At(i int) *RefinedVar {
	return t.vars[i]
}
