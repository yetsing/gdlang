package evaluator

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"weilang/ast"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

var modules = make(map[string]*object.Module)
var moduleKey = "module"

func evalExport(ctx context.Context, env *object.Environment, idents []*ast.Identifier) object.Object {
	value := ctx.Value(moduleKey)
	mod, ok := value.(*object.Module)
	if !ok {
		return object.Unreachable("not found module from context")
	}
	for _, ident := range idents {
		name := ident.Value
		if _, ok := env.Get(name); !ok {
			return object.NewError("undefined '%s'", name)
		}
		mod.AddExport(name)
	}
	return object.NULL
}

//goland:noinspection GoUnusedParameter
func evalImport(ctx context.Context, filename object.Object) object.Object {
	strObj, ok := filename.(*object.String)
	if !ok {
		return object.WrongArgumentTypeAt(filename.Type(), 1)
	}
	return importFromFile(strObj.Value)
}

func importFromFile(filename string) object.Object {
	weiFilename := filename
	if !strings.HasSuffix(filename, ".wei") {
		weiFilename = filename + ".wei"
	}
	weiFilename, _ = filepath.Abs(weiFilename)
	if _, err := os.Stat(weiFilename); errors.Is(err, os.ErrNotExist) {
		return object.NewError("Not found module filename: %s", filename)
	}

	if mod, ok := modules[weiFilename]; ok {
		return mod
	}

	l := lexer.NewWithFilename(weiFilename)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return object.NewError("%v", err)
	}

	module := object.NewModule(weiFilename)
	modules[weiFilename] = module
	evaluated := Eval(
		NewModuleContext(module),
		program,
		module.GetEnv(),
	)
	if IsError(evaluated) {
		return evaluated
	}
	return module
}

func NewModuleContext(m *object.Module) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, moduleKey, m)
}
