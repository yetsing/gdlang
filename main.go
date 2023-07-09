package main

import (
	"fmt"
	"os"
	"os/user"
	"weilang/interpreter"
	"weilang/repl"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Invalid argument")
		fmt.Println("Usage:")
		fmt.Println("    weilang [filename]")
		os.Exit(1)
	}
	// 执行文件
	if len(os.Args) == 2 {
		filename := os.Args[1]
		interpreter.RunFile(filename)
		return
	}
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Weilang programming language!\n",
		u.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
