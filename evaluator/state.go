package evaluator

import (
	"fmt"
	"os"
	"strings"
	"weilang/ast"
	"weilang/object"
)

type WeiState struct {
	module *object.Module
	stack  *object.CallStack
	// excStack 错误栈
	excStack *object.CallStack
	exc      *object.Error
}

func NewWeiState(module *object.Module) *WeiState {
	return &WeiState{
		module:   module,
		stack:    object.NewCallStack(),
		excStack: nil,
		exc:      nil,
	}
}

func (g *WeiState) CreateFrame(filename string, funcName string) *object.Frame {
	return g.stack.CreateFrame(filename, funcName)
}

func (g *WeiState) DestroyFrame() *object.Frame {
	return g.stack.DestroyFrame()
}

func (g *WeiState) Frame() *object.Frame {
	return g.stack.Top()
}

func (g *WeiState) GetModule() *object.Module {
	return g.module
}

func (g *WeiState) SetModule(module *object.Module) {
	g.module = module
}

func (g *WeiState) UpdateLocation(node ast.Node) {
	location := node.GetFileLocation()
	frame := g.Frame()
	frame.SetLineno(location.Lineno)

}

func (g *WeiState) HandleError(obj object.Object) {
	if g.exc == nil {
		g.exc = obj.(*object.Error)
		g.excStack = g.stack.Copy()
	}
}

func (g *WeiState) GetExcFrames() []*object.Frame {
	if g.excStack == nil {
		return []*object.Frame{}
	}
	return g.excStack.GetFrames()
}

func (g *WeiState) NewError(format string, args ...any) *object.Error {
	e := object.NewError(format, args...)
	g.HandleError(e)
	return e
}

func (g *WeiState) Unreachable(msg string) *object.Error {
	e := object.Unreachable(msg)
	g.HandleError(e)
	return e
}

func (g *WeiState) WrongNumberArgument(name string, got, want int) *object.Error {
	e := object.WrongNumberArgument3(name, got, want)
	g.HandleError(e)
	return e
}

func getLine(filename string, lineno int) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if lineno >= len(lines) {
		return ""
	}
	line := lines[lineno]
	return strings.TrimSpace(line)
}

func (g *WeiState) HasExc() bool {
	return g.exc != nil
}

func (g *WeiState) PrintExc() {
	// 打印错误栈
	fmt.Println("Traceback")
	for _, frame := range g.GetExcFrames() {
		fmt.Printf("  File \"%s\", line %d, in %s\n", frame.GetFilename(), frame.GetLineno()+1, frame.GetFuncName())
		fmt.Printf("    %s\n", getLine(frame.GetFilename(), frame.GetLineno()))
	}
	fmt.Println(g.exc.String())
}
