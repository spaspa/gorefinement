package refinement

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"strings"
)

func ParseWithBaseType(s string, baseType types.Type) (identName string, resultType types.Type, err error) {
	return parse(s, baseType, ":")
}

func ParseAliasWithBaseType(s string, baseType types.Type) (identName string, resultType types.Type, err error) {
	return parse(s, baseType, "=")
}

func parse(s string, baseType types.Type, sep string) (string, types.Type, error) {
	decl := strings.SplitN(s, sep, 2)

	if len(decl) != 2 {
		return "", nil, errors.New("failed to parse declaration")
	}
	resultType, err := parseType(strings.Trim(decl[1], " "), baseType)
	return decl[0], resultType, err
}

func parseType(s string, baseType types.Type) (types.Type, error) {
	typeDecl := strings.SplitN(s, "->", 2)
	if len(typeDecl) == 0 {
		return nil, errors.New("failed to parse type")
	}
	if len(typeDecl) == 2 {
		signature, ok := baseType.(*types.Signature)
		if !ok {
			return nil, errors.New("invalid type for function")
		}

		params, err := parseParamList(strings.Trim(typeDecl[0], " "), signature.Params(), false)
		if err != nil {
			return nil, err
		}
		results, err := parseParamList(strings.Trim(typeDecl[1], " "), signature.Results(), true)
		if err != nil {
			return nil, err
		}

		// TODO: add check for base signature and dependent version of signature
		return NewDependentSignature(signature, params, results), nil
	}
	return parseRefinementType(typeDecl[0], baseType)
}

func parseParamList(s string, baseTuple *types.Tuple, allowSingleParamWithoutBrace bool) (*RefinedTuple, error) {
	if !allowSingleParamWithoutBrace && (!strings.HasPrefix(s, "(") || !strings.HasSuffix(s, ")")) {
		return nil, errors.New("incorrect param list form")
	}
	body := strings.Trim(s, "() ")
	paramStrs := strings.Split(body, ",")

	numParams := len(paramStrs)
	if len(body) == 0 {
		numParams = 0
	}

	if numParams != baseTuple.Len() {
		return nil, errors.New("param list size does not match")
	}

	var results []*RefinedVar

	for i := 0; i < numParams; i++ {
		result, err := parseParam(strings.Trim(paramStrs[i], " "), baseTuple.At(i))
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return NewRefinedTuple(results...), nil
}

func parseParam(s string, baseVar *types.Var) (*RefinedVar, error) {
	var nameStr string
	var typeStr string

	// param list with single param

	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		typeStr = s
	} else {
		nameType := strings.SplitN(s, " ", 2)
		if len(nameType) == 1 {
			nameOrTypeStr := strings.Trim(nameType[0], " ")
			if strings.HasPrefix(nameOrTypeStr, "{") {
				typeStr = nameOrTypeStr
			} else {
				nameStr = nameOrTypeStr
			}
		} else if len(nameType) == 2 {
			nameStr = strings.Trim(nameType[0], " ")
			typeStr = strings.Trim(nameType[1], " ")
		}
	}

	var resultName string
	var resultType *RefinedType

	if nameStr != "" {
		expr, err := parser.ParseExpr(nameStr)
		if err != nil {
			return nil, errors.New("failed to parse param name")
		}
		ident, ok := expr.(*ast.Ident)
		if !ok {
			return nil, errors.New("failed to parse param name")
		}
		resultName = ident.Name
	}
	if typeStr != "" {
		typ, err := parseRefinementType(typeStr, baseVar.Type())
		if err != nil {
			return nil, err
		}
		resultType = typ
	}

	return &RefinedVar{
		name:        resultName,
		refinedType: resultType,
	}, nil
}

// parseRefinementType parses refinement type in form "{" var ":" type "|" predicate "}"
func parseRefinementType(s string, baseType types.Type) (*RefinedType, error) {
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") {
		return nil, errors.New("incorrect refinement form")
	}
	body := strings.Trim(s, "{} ")

	firstColonIdx := strings.Index(body, ":")
	firstVertBarIdx := strings.Index(body, "|")

	if firstVertBarIdx == -1 || firstColonIdx == -1 || firstColonIdx > firstVertBarIdx {
		return nil, errors.New("incorrect refinement form")
	}

	refinementVarStr := strings.Trim(body[:firstColonIdx], " ")
	parsedBaseTypeStr := strings.Trim(body[firstColonIdx+1:firstVertBarIdx], " ")
	predicateStr := strings.Trim(body[firstVertBarIdx+1:], " ")

	refinementVar, err := parser.ParseExpr(refinementVarStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refinement variable: %v", refinementVarStr)
	}

	refinementIdent, ok := refinementVar.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("failed to parse refinement variable: %v", refinementVarStr)
	}

	if types.TypeString(baseType, nil) != parsedBaseTypeStr {
		return nil, fmt.Errorf("base type %s is incositent inconsistent with types in go", parsedBaseTypeStr)
	}

	predicate, err := parser.ParseExpr(predicateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refinement preficate: %v", predicateStr)
	}

	return &RefinedType{
		Refinement: Refinement{
			Predicate: predicate,
			RefVar:    refinementIdent,
		},
		Type: baseType,
	}, nil
}
