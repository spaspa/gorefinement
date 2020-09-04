package liquid

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strconv"

	"github.com/spaspa/gorefinement/z3Util"

	"github.com/aclements/go-z3/z3"
)

const (
	doubleEbits = 11
	doubleSbits = 53
)

func convertToZ3Ast(env *Environment, ctx *z3.Context, expr ast.Expr) (z3.Value, error) {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExpr(env, ctx, expr)
	case *ast.UnaryExpr:
		return convertUnaryExpr(env, ctx, expr)
	case *ast.Ident:
		return convertIdent(env, ctx, expr)
	case *ast.BasicLit:
		return convertBasicLit(ctx, expr)
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported expr")
	}
}

func convertIdent(env *Environment, ctx *z3.Context, expr *ast.Ident) (z3.Value, error) {
	if expr.Name == PredicateVariableName {
		// reserved predicate variable name found
		// TODO: support non-int type
		return ctx.IntConst(PredicateVariableName), nil
	}
	_, obj := env.Scope.LookupParent(expr.Name, env.Pos)
	if obj == nil || obj.Type() == nil {
		return nil, fmt.Errorf("failed to convert expr to z3 ast: ident not found")
	}
	if basicType, ok := obj.Type().(*types.Basic); ok {
		if basicType.Info()&types.IsInteger != 0 {
			return ctx.IntConst(expr.Name), nil
		}
		if basicType.Info()&types.IsFloat != 0 {
			return ctx.RealConst(expr.Name).ToFloat(ctx.FloatSort(doubleEbits, doubleSbits)), nil
		}
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: ident is not basic type")
}

func convertBinaryExpr(env *Environment, ctx *z3.Context, expr *ast.BinaryExpr) (z3.Value, error) {
	lhsValue, err := convertToZ3Ast(env, ctx, expr.X)
	if err != nil {
		return nil, err
	}
	rhsValue, err := convertToZ3Ast(env, ctx, expr.Y)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(lhsValue) != reflect.TypeOf(rhsValue) {
		return nil, fmt.Errorf("failed to convert expr to z3 ast: type mismatch")
	}

	switch lhsValue := lhsValue.(type) {
	case z3.Int:
		return z3Util.ConvertIntBinaryExpr(lhsValue, rhsValue.(z3.Int), expr.Op)
	case z3.Bool:
		return z3Util.ConvertBoolBinaryExpr(lhsValue, rhsValue.(z3.Bool), expr.Op)
	case z3.Float:
		return z3Util.ConvertFloatBinaryExpr(lhsValue, rhsValue.(z3.Float), expr.Op)
	case z3.Real:
		return z3Util.ConvertRealBinaryExpr(lhsValue, rhsValue.(z3.Real), expr.Op)
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: not a supported type")
	}
}

func convertUnaryExpr(env *Environment, ctx *z3.Context, expr *ast.UnaryExpr) (z3.Value, error) {
	lhsValue, err := convertToZ3Ast(env, ctx, expr.X)
	if err != nil {
		return nil, err
	}

	switch lhsValue := lhsValue.(type) {
	case z3.Int:
		return z3Util.ConvertIntUnaryExpr(lhsValue, expr.Op)
	case z3.Bool:
		return z3Util.ConvertBoolUnaryExpr(lhsValue, expr.Op)
	case z3.Float:
		return z3Util.ConvertFloatUnaryExpr(lhsValue, expr.Op)
	case z3.Real:
		return z3Util.ConvertRealUnaryExpr(lhsValue, expr.Op)
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: not a supported type")
	}
}

func convertBasicLit(ctx *z3.Context, expr *ast.BasicLit) (z3.Value, error) {
	switch expr.Kind {
	case token.INT:
		if v, err := strconv.ParseInt(expr.Value, 10, 64); err == nil {
			return ctx.FromInt(v, ctx.IntSort()), nil
		}
		return nil, fmt.Errorf("failed to parse int")
	case token.FLOAT:
		if v, err := strconv.ParseFloat(expr.Value, 10); err == nil {
			// IEEE 754 double
			return ctx.FromFloat64(v, ctx.FloatSort(doubleEbits, doubleSbits)), nil
		}
		return nil, fmt.Errorf("failed to parse int")
	default:
		return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported literal")
	}
}
