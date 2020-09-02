package refinement

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("can parse simple refinement type declaration", func(t *testing.T) {
		sig := "a: { x: int | x > 0|| x== 1 }"
		typ, err := ParseWithBaseType(sig, types.Typ[types.Int])
		if err != nil {
			t.Errorf("could not parse correct refinement: %v", sig)
			t.Errorf("%s", err)
			return
		}
		refined, ok := typ.(*RefinedType)
		if !ok {
			t.Errorf("parsed to non-refined type: %v", sig)
			return
		}
		expr1, ok := refined.Predicate.(*ast.BinaryExpr)
		if !ok {
			t.Errorf("predicate was non binary expr, excepted logical or, got %v", refined.Predicate)
		}
		if expr1.Op != token.LOR {
			t.Errorf("predicate was non binary expr, excepted logical or, got %v", refined.Predicate)
		}
		fmt.Println(refined)
	})

	t.Run("can parse simple refinement type declaration", func(t *testing.T) {
		sig := "oneTwo: () -> ({ r: int | r == 1 }, { r: int | r == 2 })"
		base := `func () (int, int) { return 1, 2 }`

		baseExpr, _ := parser.ParseExpr(base)
		info := types.Info{Types: map[ast.Expr]types.TypeAndValue{}}
		_ = types.CheckExpr(token.NewFileSet(), nil, token.NoPos, baseExpr, &info)

		baseType := info.Types[baseExpr].Type
		typ, err := ParseWithBaseType(sig, baseType)
		if err != nil {
			t.Errorf("could not parse correct refinement: %v", sig)
			t.Errorf("%s", err)
			return
		}
		dependent, ok := typ.(*DependentSignature)
		if !ok {
			t.Errorf("parsed to invalid type: %v", sig)
			return
		}
		fmt.Println(dependent)
	})

	t.Run("can parse complex refinement type declaration", func(t *testing.T) {
		sig := "maxDiv: (x {v:int|true}, y, z { v: int | v > 0 }) -> { r: int | r > x / y && r > y / z }"
		base := `func (x, y, z int) int { if x / y > y / z { return x / y } else { return y / z } }`

		baseExpr, _ := parser.ParseExpr(base)
		info := types.Info{Types: map[ast.Expr]types.TypeAndValue{}}
		_ = types.CheckExpr(token.NewFileSet(), nil, token.NoPos, baseExpr, &info)

		baseType := info.Types[baseExpr].Type
		typ, err := ParseWithBaseType(sig, baseType)
		if err != nil {
			t.Errorf("could not parse correct refinement: %v", sig)
			t.Errorf("%s", err)
			return
		}
		dependent, ok := typ.(*DependentSignature)
		if !ok {
			t.Errorf("parsed to invalid type: %v", sig)
			return
		}
		fmt.Println(dependent)
	})
}
