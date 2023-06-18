# weilang

使用 Go 语言实现一门自己的解释语言

# 参考

[Let's Build A Simple Interpreter](https://github.com/rspivak/lsbasi)

《用Go语言自制解释器》

... 等网上各种文章

# 语法规范

- 基本数据类型

```text
Integer 整数，如 1 234 0b0101 0x1234DF

String 字符串，单双引号都行，如 "wei" 'wei' ，多行字符串 `abc`

Bool 布尔值， true false

null 空值
```

- 定义变量

```text
var a = 123
var b = "wei"
var c = true
var d = null
```

- 定义常量

```text
con a = 4
```

- 逻辑运算符

```text
1 == 2
1 != 2
1 > 2
1 >= 2
1 < 2
1 <= 2

not 1 == 2 and 1 == 2 or 1 == 2
```

- 算术运算符

```text
1 + 2
1 - 2
1 * 3
4 / 2
4 % 2
-1
```

- 注释

```text
// 开头到行尾都是注释
```

- 控制流

if

```text
if(conditionExpression) {
    statement1
} else if(conditionExpression) {
    statement2
} else {
    statement3
}
```

while

```text
while(conditionExpression) {
    statement
    // continue 跳过当前循环的余下逻辑, 进入下一轮循环
    continue
    // break 跳出当前循环
    break
}
```

- 函数相关

函数定义

```text
con funcName = function(para1, para2) {
    statement
    return returnValue
    // 只有 return 表示返回 null
    return
}
```

函数调用

```text
funcName(para1, para2)
```
