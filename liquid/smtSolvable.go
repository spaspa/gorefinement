package liquid

import (
	"github.com/mitchellh/go-z3"
	"log"

	"go/ast"
)

type Result int

const (
	Valid Result = iota
	Invalid
	Unknown
)

func Query(env *Environment, expr ast.Expr) Result {
	// context
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	err := config.Close()
	if err != nil {
		log.Println("z3 config close failed", err)
	}
	defer ctx.Close()

	// solver
	s := ctx.NewSolver()
	defer s.Close()

	return Unknown
}