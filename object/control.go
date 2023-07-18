package object

import "fmt"

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) TypeIs(objectType ObjectType) bool {
	return rv.Type() == objectType
}

func (rv *ReturnValue) TypeNotIs(objectType ObjectType) bool {
	return rv.Type() != objectType
}

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

type TwoReturnValue struct {
	First  Object
	Second Object
}

func (rv *TwoReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *TwoReturnValue) TypeIs(objectType ObjectType) bool {
	return rv.Type() == objectType
}

func (rv *TwoReturnValue) TypeNotIs(objectType ObjectType) bool {
	return rv.Type() != objectType
}

func (rv *TwoReturnValue) String() string {
	return fmt.Sprintf("(%s, %s)", rv.First, rv.Second)
}

// ==========================
// 两个特殊值，用于处理 continue break 语句
// ==========================
var (
	CONTINUE_VALUE = &ContinueValue{}
	BREAK_VALUE    = &BreakValue{}
)

type ContinueValue struct {
}

func (c *ContinueValue) Type() ObjectType {
	return CONTINUE_VALUE_OBJ
}

func (c *ContinueValue) TypeIs(objectType ObjectType) bool {
	return c.Type() == objectType
}

func (c *ContinueValue) TypeNotIs(objectType ObjectType) bool {
	return c.Type() != objectType
}

func (c *ContinueValue) String() string {
	return "continue"
}

type BreakValue struct {
}

func (b *BreakValue) Type() ObjectType {
	return BREAK_VALUE_OBJ
}

func (b *BreakValue) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BreakValue) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BreakValue) String() string {
	return "break"
}
