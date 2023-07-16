package evaluator

import (
	"errors"
	"os"
	"path/filepath"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

type Wei struct {
	module *object.Module
	store  map[string]object.Object
}

func (w *Wei) Type() object.ObjectType {
	return object.WEI_OBJ
}

func (w *Wei) TypeIs(objectType object.ObjectType) bool {
	return w.Type() == objectType
}

func (w *Wei) TypeNotIs(objectType object.ObjectType) bool {
	return w.Type() != objectType
}

func (w *Wei) String() string {
	return "wei"
}

func (w *Wei) SetAttribute(_ string, _ object.Object) object.Object {
	return object.NewError("'wei' do not support assignment")
}

func (w *Wei) GetAttribute(name string) object.Object {
	if value, ok := w.store[name]; ok {
		return value
	}
	return object.NewError("'wei' has not attribute '%s'", name)
}

func (w *Wei) GetModule() *object.Module {
	return w.module
}

func NewWei(filename string) *Wei {
	path := filename
	if len(filename) > 0 {
		path, _ = filepath.Abs(filename)
	}
	store := make(map[string]object.Object)
	w := &Wei{module: object.NewModule(path), store: store}
	store["filename"] = object.NewString(path)
	for s, builtin := range weiBuiltins {
		store[s] = builtin
	}
	env := w.GetModule().GetEnv()
	env.Add("wei", w, true)
	return w
}

func NewWeiEnvironment(filename string) *object.Environment {
	w := NewWei(filename)
	return w.module.GetEnv()
}

var weiBuiltins map[string]*object.Builtin

func init() {
	weiBuiltins = map[string]*object.Builtin{
		"import": {
			Name: "import",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.WrongNumberArgument(len(args), 1)
				}
				arg, ok := args[0].(*object.String)
				if !ok {
					return object.WrongArgumentTypeAt(args[0].Type(), 1)
				}
				filename := arg.Value
				path, _ := filepath.Abs(filename)
				if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
					return object.NewError("Not found module filename: %s", filename)
				}
				l := lexer.NewWithFilename(path)
				p := parser.New(l)
				program, err := p.ParseProgram()
				if err != nil {
					return object.NewError("%v", err)
				}

				wei := NewWei(path)
				env := wei.module.GetEnv()
				env.Add("wei", wei, true)
				evaluated := Eval(program, env)
				if IsError(evaluated) {
					return evaluated
				}
				return wei.module
			},
		},
	}
}
