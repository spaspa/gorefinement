package liquid

import (
	"fmt"
	"github.com/spaspa/gorefinement/refinement"
	"go/types"
)

// IsSubType checks t1 <: t2 or not.
func IsSubtype(e *Environment, t1, t2 types.Type) bool {
	switch t1 := t1.(type) {
	case *refinement.RefinedType:
		return isSubTypeOfRefinedType(e, t1, t2)
	case *refinement.DependentSignature:
		panic("not supported")
	default:
		return isSubTypeOfNonRefinedType(e, t1, t2)
	}
}

// isSubTypeOfRefinedType checks t1 <: t2 or not.
// t1 should be refined type.
func isSubTypeOfRefinedType(e *Environment, t1 *refinement.RefinedType, t2 types.Type) bool {
	switch t2 := t2.(type) {
	case *refinement.RefinedType:
		return verify(e, t1, t2)
	case *refinement.DependentSignature:
		panic("not supported")
	default:
		return types.Identical(t1.Type, t2)
	}
}

// isSubTypeOfNonRefinedType checks t1 <: t2 or not.
// t1 should be non-refined type.
func isSubTypeOfNonRefinedType(e *Environment, t1, t2 types.Type) bool {
	switch t2 := t2.(type) {
	case *refinement.RefinedType:
		return types.Identical(t1, t2.Type) && verify(e, refinement.NewRefinedTypeWithTruePredicate(t1), t2)
	case *refinement.DependentSignature:
		panic("not supported")
	default:
		return types.Identical(t1, t2)
	}
}

// verify checks implication ⟦e⟧ ∧ ⟦r1⟧ ⇒ ⟦r2⟧ is valid or not.
func verify(e *Environment, r1, r2 *refinement.RefinedType) bool {
	// TODO implement
	fmt.Println(types.ExprString(e.Embedding()))
	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println("-----")
	return true
}
