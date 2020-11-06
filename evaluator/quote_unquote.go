package evaluator

import (
	"github.com/tatsuya4559/monkey/ast"
	"github.com/tatsuya4559/monkey/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
