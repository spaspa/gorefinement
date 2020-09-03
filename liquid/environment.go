package liquid

import "go/types"

type ObjectRefinementMap = map[types.Object]types.Type

type Environment struct {
	ExplicitRefinementMap ObjectRefinementMap
	ImplicitRefinementMap ObjectRefinementMap
	FunArgRefinementMap   ObjectRefinementMap
	Scope *types.Scope
}

func NewEnvironment() *Environment {
	return &Environment{
		ExplicitRefinementMap: map[types.Object]types.Type{},
		ImplicitRefinementMap: map[types.Object]types.Type{},
		FunArgRefinementMap:   map[types.Object]types.Type{},
		Scope:                 nil,
	}
}

func (env *Environment)RefinementTypeOf(object types.Object) types.Type {
	e := env.ExplicitRefinementMap[object]
	if e != nil {
		return e
	}
	i := env.ImplicitRefinementMap[object]
	if i != nil {
		return i
	}
	return nil
}
