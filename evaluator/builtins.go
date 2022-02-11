package evaluator

import (
	"fmt"

	"github.com/tatsuya4559/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: _len},
	"first": {Fn: _first},
	"last":  {Fn: _last},
	"rest":  {Fn: _rest},
	"push":  {Fn: _push},
	"puts":  {Fn: _puts},
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. want=1, got=%d", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Hash:
		return &object.Integer{Value: int64(len(arg.Pairs))}
	default:
		return newError("argument to `len` not supported, got %s",
			arg.Type())
	}
}

func _first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. want=1, got=%d", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func _last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. want=1, got=%d", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}
	return NULL
}

func _rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. want=1, got=%d", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}
	return NULL
}

func _push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. want=2, got=%d", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("first argument to `push` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func _puts(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}
