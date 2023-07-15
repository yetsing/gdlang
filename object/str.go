package object

import (
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"hash/fnv"
	"strings"
	"unicode/utf8"
)

type String struct {
	*attributeStore
	Length int
	Value  string
}

func NewString(val string) *String {
	return &String{
		attributeStore: strAttr,
		Length:         utf8.RuneCountInString(val),
		Value:          val,
	}
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) TypeIs(objectType ObjectType) bool {
	return s.Type() == objectType
}

func (s *String) TypeNotIs(objectType ObjectType) bool {
	return s.Type() != objectType
}

func (s *String) String() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

func (s *String) GetAttribute(name string) Object {
	ret := s.attributeStore.get(s, name)
	if ret != nil {
		return ret
	}
	return attributeError(string(s.Type()), name)
}

func (s *String) SetAttribute(name string, _ Object) Object {
	return attributeError(string(s.Type()), name)
}

func (s *String) slice(start, end int) string {
	start = convertRange(start, s.Length)
	end = convertRange(end, s.Length)
	return exutf8.RuneSubString(s.Value, start, end-start)
}

// ================================
// str 对象的内置属性和方法
// ================================
var strAttr *attributeStore

// countMethod 计算子字符串出现次数
//
// args:
//
//	sub: str 必填 子字符串
//	start: int 选填 开始位置
//	end: int 选填 结束位置
//
// 有这么几种调用方式
// countMethod(s, sub)
// countMethod(s, sub, start)
// countMethod(s, sub, start, end)
func countMethod(obj Object, args ...Object) Object {
	argc := len(args)
	if argc < 1 || argc > 3 {
		return WrongNumberArgument2(argc, 1, 3)
	}

	this := obj.(*String)
	sub, ok := args[0].(*String)
	if !ok {
		return wrongArgumentTypeAt(args[0].Type(), 1)
	}
	if argc == 1 {
		return NewInteger(int64(strings.Count(this.Value, sub.Value)))
	}

	startObj, ok := args[1].(*Integer)
	if !ok {
		return wrongArgumentTypeAt(args[1].Type(), 2)
	}
	start := int(startObj.Value)
	end := this.Length
	if argc == 3 {
		endObj, ok := args[2].(*Integer)
		if !ok {
			return wrongArgumentTypeAt(args[2].Type(), 3)
		}
		end = int(endObj.Value)
	}

	if end-start < sub.Length {
		return NewInteger(0)
	}
	s := this.slice(start, end)
	n := strings.Count(s, sub.Value)
	return NewInteger(int64(n))
}

// endswithMethod 字符串是否有指定后缀
// str.endswith(suffix[, start[, end]])
func endswithMethod(obj Object, args ...Object) Object {
	argc := len(args)
	if argc < 1 || argc > 3 {
		return WrongNumberArgument2(argc, 1, 3)
	}

	this := obj.(*String)
	start := 0
	end := this.Length

	sub, ok := args[0].(*String)
	if !ok {
		return wrongArgumentTypeAt(args[0].Type(), 1)
	}
	if argc > 1 {
		startObj, ok := args[1].(*Integer)
		if !ok {
			return wrongArgumentTypeAt(args[1].Type(), 2)
		}
		start = int(startObj.Value)
		start = convertRange(start, this.Length)
		if argc == 3 {
			endObj, ok := args[2].(*Integer)
			if !ok {
				return wrongArgumentTypeAt(args[2].Type(), 3)
			}
			end = int(endObj.Value)
			end = convertRange(end, this.Length)
		}
	}

	if sub.Length == 0 {
		return NativeBoolToBooleanObject(true)
	}

	if end-start < sub.Length {
		return NativeBoolToBooleanObject(false)
	}

	s := this.slice(end-sub.Length, end)
	return NativeBoolToBooleanObject(s == sub.Value)
}

// findMethod 返回子字符串第一次出现的位置
// str.find(sub[, start[, end]])
func findMethod(obj Object, args ...Object) Object {
	argc := len(args)
	if argc < 1 || argc > 3 {
		return WrongNumberArgument2(argc, 1, 3)
	}
	this := obj.(*String)

	sub, ok := args[0].(*String)
	if !ok {
		return wrongArgumentTypeAt(args[0].Type(), 1)
	}

	if argc == 1 {
		byteIndex := strings.Index(this.Value, sub.Value)
		if byteIndex != -1 {
			byteIndex = utf8.RuneCountInString(this.Value[:byteIndex])
		}
		return NewInteger(int64(byteIndex))
	}

	startObj, ok := args[1].(*Integer)
	if !ok {
		return wrongArgumentTypeAt(args[1].Type(), 2)
	}
	start := int(startObj.Value)
	end := this.Length
	if argc == 3 {
		endObj, ok := args[2].(*Integer)
		if !ok {
			return wrongArgumentTypeAt(args[2].Type(), 3)
		}
		end = int(endObj.Value)
	}
	s := this.slice(start, end)
	byteIndex := strings.Index(s, sub.Value)
	if byteIndex != -1 {
		byteIndex = utf8.RuneCountInString(s[:byteIndex]) + convertRange(start, this.Length)
	}
	return NewInteger(int64(byteIndex))
}

// formatMethod 格式化字符串
// format(*args)
// 使用的占位符为 {} ，只有这一种
// 如果要输入原始的 '{' '}' 符号，使用 '{{' 表示 '{' ，'}}' 表示 '}'
//
// 例子
//
//	'a {}'.format(1) => 'a 1'
//	'a {{}}'.format() => 'a {}'
func formatMethod(obj Object, args ...Object) Object {
	this := obj.(*String)
	byteCount := len(this.Value)
	want := strings.Count(this.Value, "{}")
	argc := len(args)
	if argc != want {
		return WrongNumberArgument(argc, want)
	}
	var out strings.Builder
	out.Grow(byteCount)
	argIndex := 0
	var found byte
	for i := 0; i < byteCount; i++ {
		c := this.Value[i]
		if found != 0 {
			// 两个相同的花括号， '{{' 或者 '}}'
			if c == found {
				out.WriteByte(c)
				found = 0
			} else if c == '}' {
				out.WriteString(args[argIndex].String())
				argIndex++
				found = 0
			} else {
				return NewError("single '%s' encountered in format string", string(found))
			}
			continue
		}

		if c == '{' || c == '}' {
			found = c
		} else {
			out.WriteByte(c)
		}
	}
	if found != 0 {
		return NewError("single '%s' encountered in format string", string(found))
	}
	if argIndex != argc {
		return WrongNumberArgument(argc, argIndex)
	}
	return NewString(out.String())
}

// joinMethod 连接数组中的对象字符串，分隔符为调用的字符串
// str.join(list)
//
// 例子
//
//	','.join([1, 2, 3]) => '1,2,3'
func joinMethod(obj Object, args ...Object) Object {
	if len(args) != 1 {
		return WrongNumberArgument(len(args), 1)
	}
	this := obj.(*String)
	arg, ok := args[0].(*List)
	if !ok {
		return wrongArgumentType(args[0].Type())
	}
	var objs []string
	for _, element := range arg.Elements {
		objs = append(objs, element.String())
	}
	s := strings.Join(objs, this.Value)
	return NewString(s)
}

// lowerMethod 字符串转小写
// str.lower()
func lowerMethod(obj Object, args ...Object) Object {
	if len(args) > 0 {
		return WrongNumberArgument(len(args), 0)
	}
	this := obj.(*String)
	s := strings.ToLower(this.Value)
	return NewString(s)
}

// splitMethod 字符串分割
// str.split(sep, [maxsplit])
//
// 例子
//
//	"a,b,c".split(",") => ["a", "b", "c"]
//	"a,b,c".split(",", 0) => ["a,b,c"]
//	"a,b,c".split(",", 1) => ["a", "b,c"]
func splitMethod(obj Object, args ...Object) Object {
	argc := len(args)
	if argc < 1 || argc > 2 {
		return WrongNumberArgument2(argc, 1, 2)
	}

	this := obj.(*String)
	sepObj, ok := args[0].(*String)
	if !ok {
		return wrongArgumentTypeAt(args[0].Type(), 1)
	}
	if sepObj.Length == 0 {
		return NewError("empty separator")
	}
	if argc == 1 {
		result := strings.Split(this.Value, sepObj.Value)
		var elements []Object
		for _, s := range result {
			elements = append(elements, NewString(s))
		}
		return NewList(elements)
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		return wrongArgumentTypeAt(args[1].Type(), 2)
	}
	maxsplit := int(intObj.Value)
	result := strings.SplitN(this.Value, sepObj.Value, maxsplit+1)
	var elements []Object
	for _, s := range result {
		elements = append(elements, NewString(s))
	}
	return NewList(elements)
}

// startswithMethod 如果字符串以 prefix 开头，返回 true ；否则返回 false
// str.startswith(prefix[, start[, end]])
func startswithMethod(obj Object, args ...Object) Object {
	argc := len(args)
	if argc < 1 || argc > 3 {
		return WrongNumberArgument2(argc, 1, 3)
	}

	this := obj.(*String)
	start := 0
	end := this.Length

	sub, ok := args[0].(*String)
	if !ok {
		return wrongArgumentTypeAt(args[0].Type(), 1)
	}
	if argc > 1 {
		startObj, ok := args[1].(*Integer)
		if !ok {
			return wrongArgumentTypeAt(args[1].Type(), 2)
		}
		start = convertRange(int(startObj.Value), this.Length)
		if argc == 3 {
			endObj, ok := args[2].(*Integer)
			if !ok {
				return wrongArgumentTypeAt(args[2].Type(), 3)
			}
			end = convertRange(int(endObj.Value), this.Length)
		}
	}

	if sub.Length == 0 {
		return NativeBoolToBooleanObject(true)
	}

	if end-start < sub.Length {
		return NativeBoolToBooleanObject(false)
	}

	s := this.slice(start, start+sub.Length)
	return NativeBoolToBooleanObject(s == sub.Value)
}

// stripMethod 移除字符串前后指定字符
// str.strip(chars)
func stripMethod(obj Object, args ...Object) Object {
	if len(args) != 1 {
		return WrongNumberArgument(len(args), 1)
	}
	this := obj.(*String)
	arg, ok := args[0].(*String)
	if !ok {
		return wrongArgumentType(args[0].Type())
	}
	s := strings.Trim(this.Value, arg.Value)
	return NewString(s)
}

func upperMethod(obj Object, args ...Object) Object {
	if len(args) != 0 {
		return WrongNumberArgument(len(args), 0)
	}
	this := obj.(*String)
	return NewString(strings.ToUpper(this.Value))
}

func init() {
	strAttr = &attributeStore{
		attribute: map[string]Object{
			"count": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "count",
				Fn:    countMethod,
			},
			"endswith": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "endswith",
				Fn:    endswithMethod,
			},
			"find": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "find",
				Fn:    findMethod,
			},
			"format": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "format",
				Fn:    formatMethod,
			},
			"join": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "join",
				Fn:    joinMethod,
			},
			"lower": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "lower",
				Fn:    lowerMethod,
			},
			"split": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "split",
				Fn:    splitMethod,
			},
			"startswith": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "startswith",
				Fn:    startswithMethod,
			},
			"strip": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "strip",
				Fn:    stripMethod,
			},
			"upper": &BuiltinMethod{
				ctype: STRING_OBJ,
				name:  "upper",
				Fn:    upperMethod,
			},
		},
	}
}
