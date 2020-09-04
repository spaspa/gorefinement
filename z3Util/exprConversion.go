package z3Util

import (
	"fmt"
	"go/token"

	"github.com/aclements/go-z3/z3"
)

func ConvertIntBinaryExpr(lhs, rhs z3.Int, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return lhs.Add(rhs), nil
	case token.SUB:
		return lhs.Sub(rhs), nil
	case token.MUL:
		return lhs.Mul(rhs), nil
	case token.QUO:
		return lhs.Div(rhs), nil
	case token.REM:
		return lhs.Rem(rhs), nil
	case token.EQL:
		return lhs.Eq(rhs), nil
	case token.LSS:
		return lhs.LT(rhs), nil
	case token.GTR:
		return lhs.GT(rhs), nil
	case token.NEQ:
		return lhs.NE(rhs), nil
	case token.LEQ:
		return lhs.LE(rhs), nil
	case token.GEQ:
		return lhs.GE(rhs), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported binary expr")
}

func ConvertFloatBinaryExpr(lhs, rhs z3.Float, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return lhs.Add(rhs), nil
	case token.SUB:
		return lhs.Sub(rhs), nil
	case token.MUL:
		return lhs.Mul(rhs), nil
	case token.QUO:
		return lhs.Div(rhs), nil
	case token.REM:
		return lhs.Rem(rhs), nil
	case token.EQL:
		return lhs.Eq(rhs), nil
	case token.LSS:
		return lhs.LT(rhs), nil
	case token.GTR:
		return lhs.GT(rhs), nil
	case token.NEQ:
		return lhs.NE(rhs), nil
	case token.LEQ:
		return lhs.LE(rhs), nil
	case token.GEQ:
		return lhs.GE(rhs), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported binary expr")
}
func ConvertRealBinaryExpr(lhs, rhs z3.Real, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return lhs.Add(rhs), nil
	case token.SUB:
		return lhs.Sub(rhs), nil
	case token.MUL:
		return lhs.Mul(rhs), nil
	case token.QUO:
		return lhs.Div(rhs), nil
	case token.EQL:
		return lhs.Eq(rhs), nil
	case token.LSS:
		return lhs.LT(rhs), nil
	case token.GTR:
		return lhs.GT(rhs), nil
	case token.NEQ:
		return lhs.NE(rhs), nil
	case token.LEQ:
		return lhs.LE(rhs), nil
	case token.GEQ:
		return lhs.GE(rhs), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported binary expr")
}

func ConvertBoolBinaryExpr(lhs, rhs z3.Bool, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.LAND:
		return lhs.And(rhs), nil
	case token.LOR:
		return lhs.Or(rhs), nil
	case token.EQL:
		return lhs.Eq(rhs), nil
	case token.NEQ:
		return lhs.NE(rhs), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported binary expr")
}

func ConvertIntUnaryExpr(arg z3.Int, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return arg, nil
	case token.SUB:
		return arg.Neg(), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported unary expr")
}

func ConvertFloatUnaryExpr(arg z3.Float, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return arg, nil
	case token.SUB:
		return arg.Neg(), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported unary expr")
}

func ConvertRealUnaryExpr(arg z3.Real, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.ADD:
		return arg, nil
	case token.SUB:
		return arg.Neg(), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported unary expr")
}

func ConvertBoolUnaryExpr(arg z3.Bool, tok token.Token) (z3.Value, error) {
	switch tok {
	case token.NOT:
		return arg.Not(), nil
	}
	return nil, fmt.Errorf("failed to convert expr to z3 ast: unsupported unary expr")
}
