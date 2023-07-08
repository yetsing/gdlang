package evaluator

import (
	"testing"
	"weilang/object"
)

func TestStringBuiltinAttributeReference(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`'abc'.ddd(1)`, "'str' object has not attribute 'ddd'", true},
		{`'abc'.w`, "'str' object has not attribute 'w'", true},

		{`'abc'.upper()`, "ABC", false},
		{`'a'.upper(); "Abc".upper()`, "ABC", false},
		{`'a中文'.upper()`, "A中文", false},
		{`'abc'.upper(1)`, "wrong number of arguments. got=1, want=0", true},

		{`'abc'.count()`, "wrong number of arguments. got=0, want=1-3", true},
		{`'abc'.count('b', 1, 2, 3)`, "wrong number of arguments. got=4, want=1-3", true},
		{`'abc'.count('a')`, 1, false},
		{`'abc'.count('a', 1)`, 0, false},
		{`'abc'.count('a', -2)`, 0, false},
		{`'abcabc'.count('a', 0)`, 2, false},
		{`'abcabc'.count('a', -6)`, 2, false},
		{`'abcabc'.count('a', 1, 2)`, 0, false},
		{`'abcabc'.count('a', -5, -3)`, 0, false},
		{`'中文'.count('a')`, 0, false},
		{`'中文'.count('文')`, 1, false},
		{`'中文'.count('文', 2)`, 0, false},
		{`'中文'.count('文', -4)`, 1, false},
		{`'中文中文'.count('文', 1, 2)`, 1, false},
		{`'中文中文'.count('文', -100, -200)`, 0, false},
		{`'中文中文'.count('中文')`, 2, false},
		{`'中文a中文'.count('a中文')`, 1, false},

		{`'中文'.endswith()`, "wrong number of arguments. got=0, want=1-3", true},
		{`'中文'.endswith('a', 1, 2, 3)`, "wrong number of arguments. got=4, want=1-3", true},
		{`'中文'.endswith(1)`, "wrong argument type: 'int' at 1", true},
		{`'中文'.endswith('', '')`, "wrong argument type: 'str' at 2", true},
		{`'中文'.endswith('')`, true, false},
		{`'中文'.endswith('a')`, false, false},
		{`'中文'.endswith('文')`, true, false},
		{`'中文'.endswith('文', 0, 1)`, false, false},
		{`'中文abc'.endswith('', 0, 1)`, true, false},
		{`'中文abc'.endswith('abc')`, true, false},
		{`'中文abc'.endswith('abc', -5)`, true, false},
		{`'中文abc'.endswith('ab', -5, -1)`, true, false},

		{`''.find()`, "wrong number of arguments. got=0, want=1-3", true},
		{`''.find(1)`, "wrong argument type: 'int' at 1", true},
		{`''.find('')`, 0, false},
		{`'abc'.find('b')`, 1, false},
		{`'abcdfddfas'.find('bc')`, 1, false},
		{`'abcdfddfas'.find('s')`, 9, false},
		{`'abcdfddfas'.find('saas')`, -1, false},
		{`'中文'.find('中文')`, 0, false},
		{`'a中文a'.find('中文')`, 1, false},
		{`'中文'.find('百度')`, -1, false},
		{`''.find('', '')`, "wrong argument type: 'str' at 2", true},
		{`''.find('', 0)`, 0, false},
		{`''.find('', 1)`, 0, false},
		{`'abc'.find('', 1)`, 1, false},
		{`'abc'.find('中', -100)`, -1, false},
		{`'abc中文'.find('中文', 1)`, 3, false},
		{`'abc中文'.find('中文', 2)`, 3, false},
		{`'abc中文'.find('中文', 3)`, 3, false},
		{`'abc中文'.find('中文', -2)`, 3, false},
		{`'中文abc中文'.find('中文', -2)`, 5, false},
		{`'中文abc中文'.find('中文', 1)`, 5, false},
		{`''.find('', 1, '')`, "wrong argument type: 'str' at 3", true},
		{`''.find('', 0, 100)`, 0, false},
		{`''.find('a', 0, 100)`, -1, false},
		{`'abc'.find('a', 0, 100)`, 0, false},
		{`'中文abc'.find('a', 0, 100)`, 2, false},
		{`'中文abc'.find('a', 2, 100)`, 2, false},
		{`'中文abc'.find('a', 2, -2)`, 2, false},
		{`'中文abc'.find('a', -4, -2)`, 2, false},
		{`'中文abc'.find('a', -100, 200)`, 2, false},
		{`'本节将给解释器添加内置函数。'.find('。', 5, 100)`, 13, false},

		{`'abc'.format()`, "abc", false},
		{`'abc'.format(1)`, "wrong number of arguments. got=1, want=0", true},
		{`'abc{{}}'.format(1)`, "wrong number of arguments. got=1, want=0", true},
		{`'abc{} {}'.format(1)`, "wrong number of arguments. got=1, want=2", true},
		{`'abc{'.format()`, "single '{' encountered in format string", true},
		{`'abc{ '.format()`, "single '{' encountered in format string", true},
		{`'abc{ }'.format()`, "single '{' encountered in format string", true},
		{`'abc{abc'.format()`, "single '{' encountered in format string", true},
		{`'abc}'.format()`, "single '}' encountered in format string", true},
		{`'abc} '.format()`, "single '}' encountered in format string", true},
		{`'hello {}'.format('john')`, "hello john", false},
		{`'hello {}'.format(123)`, "hello 123", false},
		{`'hello {{{}}}'.format(123)`, "hello {123}", false},
		{`'hello {{    }}, you are k'.format()`, "hello {    }, you are k", false},
		{`'hello {} {}'.format(1, "中文")`, "hello 1 中文", false},
		{`'hello {} {} {}'.format(1, "中文", true)`, "hello 1 中文 true", false},
		{`'hello {} {} {} }}'.format(1, "中文", true)`, "hello 1 中文 true }", false},
		{`'hello {} {} {} {{'.format(1, "中文", true)`, "hello 1 中文 true {", false},
		{`'你好 {} {} {} {{'.format(1, "中文", true)`, "你好 1 中文 true {", false},

		{`",".join()`, "wrong number of arguments. got=0, want=1", true},
		{`",".join(1, 2)`, "wrong number of arguments. got=2, want=1", true},
		{`",".join(1)`, "wrong argument type: 'int'", true},
		{`",".join('abc')`, "wrong argument type: 'str'", true},
		{`",".join({})`, "wrong argument type: 'dict'", true},
		{`",".join([])`, "", false},
		{`",".join([1, 2, 3])`, "1,2,3", false},
		{`"分隔符".join([1, 2, '3'])`, "1分隔符2分隔符3", false},
		{`"分隔符".join([1, true, null, '手掌'])`, "1分隔符true分隔符null分隔符手掌", false},

		{`"".lower(1)`, "wrong number of arguments. got=1, want=0", true},
		{`"".lower()`, "", false},
		{`"abc".lower()`, "abc", false},
		{`"Abc".lower()`, "abc", false},
		{`"Abc中文".lower()`, "abc中文", false},

		{`"".split()`, "wrong number of arguments. got=0, want=1-2", true},
		{`"".split("", 2, 2)`, "wrong number of arguments. got=3, want=1-2", true},
		{`"".split(1)`, "wrong argument type: 'int' at 1", true},
		{`"".split("")`, "empty separator", true},
		{`"".split("a", "")`, "wrong argument type: 'str' at 2", true},
		{`"a,b,c".split(",")`, []string{"a", "b", "c"}, true},
		{`"a，b，c".split("，")`, []string{"a", "b", "c"}, true},
		{`"a，b，c".split("b")`, []string{"a，", "，c"}, true},
		{`"编程语言".split("b")`, []string{"编程语言"}, true},
		{`"a，b，c".split("，", 0)`, []string{"a，b，c"}, true},
		{`"a，b，c".split("，", 1)`, []string{"a", "b，c"}, true},
		{`"a，b，c".split("，", 2)`, []string{"a", "b", "c"}, true},
		{`"a，b，c".split("，", 3)`, []string{"a", "b", "c"}, true},
		{`"中文 英文 德文".split(" ", 3)`, []string{"中文", "英文", "德文"}, true},

		{`"abc".startswith()`, "wrong number of arguments. got=0, want=1-3", true},
		{`"abc".startswith(1)`, "wrong argument type: 'int' at 1", true},
		{`"abc".startswith("")`, true, false},
		{`"abc".startswith("a")`, true, false},
		{`"abc".startswith("ab")`, true, false},
		{`"abc".startswith("abc")`, true, false},
		{`"abc".startswith("abcd")`, false, false},
		{`"abc中文".startswith("abc")`, true, false},
		{`"中文abc".startswith("中文")`, true, false},
		{`"中文abc".startswith("文a")`, false, false},
		{`"abc".startswith("", "")`, "wrong argument type: 'str' at 2", true},
		{`"abc".startswith("", 2, "")`, "wrong argument type: 'str' at 3", true},
		{`"中文abc".startswith("文a", 1)`, true, false},
		{`"中文abc".startswith("文a", 1, 2)`, false, false},
		{`"中文abc".startswith("文a", 1, -2)`, true, false},
		{`"中文abc".startswith("文a", -4, -2)`, true, false},
		{`"中文abc".startswith("abc", -3)`, true, false},
		{`"中文abc".startswith("abc", -3, -2)`, false, false},
		{`"中文abc".startswith("abc", -3, 300)`, true, false},

		{`"中文".strip()`, "wrong number of arguments. got=0, want=1", true},
		{`"中文".strip(1)`, "wrong argument type: 'int'", true},
		{`"abc".strip('', 1)`, "wrong number of arguments. got=2, want=1", true},
		{`"  abc  ".strip('')`, "  abc  ", false},
		{`"  abc  ".strip(' ')`, "abc", false},
		{`"  中文  ".strip(' ')`, "中文", false},
		{`"中文中文".strip('中文')`, "", false},
		{`"中文abc中文".strip('中文')`, "abc", false},
		{`"a b".strip(' ')`, "a b", false},
		{`'www.example.com'.strip('cmowz.')`, "example", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			if tt.isError {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. got=%T (%+v)",
						evaluated, evaluated)
					continue
				}
				if errObj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q",
						expected, errObj.Message)
				}
			} else {
				strObj, ok := evaluated.(*object.String)
				if !ok {
					t.Errorf("object is not string. got=%T (%+v)",
						evaluated, evaluated)
					continue
				}
				if strObj.Value != expected {
					t.Errorf("wrong string value. expected=%q, got=%q",
						expected, strObj.Value)
				}
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		case []string:
			list, ok := evaluated.(*object.List)
			if !ok {
				t.Errorf("object is not List. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(list.Elements) != len(expected) {
				t.Errorf("list length not equal. got=%d, want=%d", len(list.Elements), len(expected))
				continue
			}
			for i, s := range expected {
				ele := list.Elements[i]
				testStringObject(t, ele, s)
			}
		default:
			t.Errorf("invalid case")
		}
	}
}
