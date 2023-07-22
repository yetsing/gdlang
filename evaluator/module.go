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

func CacheModule(module *object.Module) {
	modules[module.Filename()] = module
}

//goland:noinspection GoUnusedParameter
func evalExport(ctx context.Context, state *WeiState, env *object.Environment, idents []*ast.Identifier) object.Object {
	mod := state.GetModule()
	for _, ident := range idents {
		state.UpdateLocation(ident)
		name := ident.Value
		if _, ok := env.Get(name); !ok {
			obj := object.NewError("undefined '%s'", name)
			state.HandleError(obj)
			return obj
		}
		mod.AddExport(name)
	}
	return nil
}

//goland:noinspection GoUnusedParameter
func evalImport(ctx context.Context, state *WeiState, filename string) object.Object {
	origModule := state.GetModule()
	ret := importFromFile(ctx, state, filename)
	state.SetModule(origModule)
	return ret
}

func importFromFile(ctx context.Context, state *WeiState, filename string) object.Object {
	weiFilename := filename
	if !strings.HasSuffix(filename, ".wei") {
		weiFilename = filename + ".wei"
	}
	weiFilename, _ = filepath.Abs(weiFilename)
	if _, err := os.Stat(weiFilename); errors.Is(err, os.ErrNotExist) {
		return state.NewError("Not found module: %s", filename)
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
	CacheModule(module)
	state.SetModule(module)
	state.CreateFrame(weiFilename, "import")
	evaluated := Eval(
		ctx,
		state,
		program,
		module.GetEnv(),
	)
	state.DestroyFrame()
	if IsError(evaluated) {
		return evaluated
	}
	return module
}
