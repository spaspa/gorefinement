package liquid

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strconv"

	"github.com/mitchellh/go-z3"
)

func ConvertToZ3Ast(env *Environment, ctx *z3.Context, expr ast.Expr) (*z3.AST, error) {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExpr(env, ctx, expr)
	case *ast.UnaryExpr:
		return convertUnaryExpr(env, ctx, expr)
	case *ast.Ident:
		return convertIdent(env, ctx, expr)
	case *ast.BasicLit:
		return convertBasicLit(env, ctx, expr)
	case nil:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: found nil expr")
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: %s is unsupported", reflect.ValueOf(expr))
	}
}

func convertIdent(env *Environment, ctx *z3.Context, expr *ast.Ident) (*z3.AST, error) {
	if expr.Name == predicateVariableName {
		// reserved predicate variable name found
		// TODO: support non-int type
		return ctx.Const(ctx.Symbol(predicateVariableName), ctx.IntSort()), nil
	}
	_, obj := env.Scope.LookupParent(expr.Name, env.Pos)
	if obj == nil || obj.Type() == nil {
		return lookupFunArgIdent(env, ctx, expr.Name)
	}
	if basicType, ok := obj.Type().(*types.Basic); ok {
		if basicType.Info()&types.IsInteger != 0 {
			return ctx.Const(ctx.Symbol(expr.Name), ctx.IntSort()), nil
		}
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: ident is not basic type")
}

func lookupFunArgIdent(env *Environment, ctx *z3.Context, name string) (*z3.AST, error) {
	if _, ok := env.FunArgRefinementMap[name]; ok {
		// TODO: support non-int type
		return ctx.Const(ctx.Symbol(name), ctx.IntSort()), nil
	}
	return nil, fmt.Errorf(`failed to convert expr to z3 ast: ident "%s" not found`, name)
}

func convertBinaryExpr(env *Environment, ctx *z3.Context, expr *ast.BinaryExpr) (*z3.AST, error) {
	lhs, err := ConvertToZ3Ast(env, ctx, expr.X)
	if err != nil {
		return nil, err
	}
	rhs, err := ConvertToZ3Ast(env, ctx, expr.Y)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case token.ADD:
		return lhs.Add(rhs), nil
	case token.LAND:
		return lhs.And(rhs), nil
	case token.EQL:
		return lhs.Eq(rhs), nil
	case token.GTR:
		return lhs.Gt(rhs), nil
	case token.GEQ:
		return lhs.Ge(rhs), nil
	case token.LEQ:
		return lhs.Le(rhs), nil
	case token.LSS:
		return lhs.Lt(rhs), nil
	case token.MUL:
		return lhs.Mul(rhs), nil
	case token.SUB:
		return lhs.Sub(rhs), nil
	case token.XOR:
		return lhs.Xor(rhs), nil
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported binary op")
	}
}

func convertUnaryExpr(env *Environment, ctx *z3.Context, expr *ast.UnaryExpr) (*z3.AST, error) {
	lhs, err := ConvertToZ3Ast(env, ctx, expr.X)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case token.SUB:
		zero := ctx.Int(0, ctx.IntSort())
		return zero.Sub(lhs), nil
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported unary op")
	}
}

func convertBasicLit(_ *Environment, ctx *z3.Context, expr *ast.BasicLit) (*z3.AST, error) {
	switch expr.Kind {
	case token.INT:
		if v, err := strconv.Atoi(expr.Value); err == nil {
			return ctx.Int(v, ctx.IntSort()), nil
		}
		return nil, fmt.Errorf("failed to parse int")
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported literal")
	}
}
